package app

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sbx/internal/process"
	"sbx/internal/tui"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/charmbracelet/x/term"
)

func (app *App) cmdNew(args []string) (returnErr error) {
	name := ""
	prompt := ""
	agent := app.config.Agent
	agentSet, promptSet := false, false
	noOpen := false
	allProjects := false
	for len(args) > 0 {
		argument := args[0]
		args = args[1:]
		switch argument {
		case "--agent", "--model", "--prompt", "--name":
			if len(args) == 0 {
				return fmt.Errorf("%s requires a value", argument)
			}
			value := args[0]
			args = args[1:]
			switch argument {
			case "--agent", "--model":
				agent, agentSet = value, true
			case "--prompt":
				prompt, promptSet = value, true
			case "--name":
				if name != "" {
					return errors.New("new accepts only one name")
				}
				name = value
			}
		case "--no-open":
			noOpen = true
		case "--all":
			allProjects = true
		case "-h", "--help":
			fmt.Fprintf(app.Out, "Usage: %s new [--all] [--no-open] [--agent claude|codex] [--prompt text] [name]\n", app.Program)
			return nil
		default:
			if strings.HasPrefix(argument, "-") {
				return fmt.Errorf("unknown option for new: %s", argument)
			}
			if name != "" {
				return errors.New("new accepts only one name")
			}
			name = argument
		}
	}
	projects, err := app.loadProjects()
	if err != nil {
		return err
	}
	selected := make(map[string]bool, len(projects))
	formUsed := false
	input, inputOK := app.In.(*os.File)
	output, outputOK := app.Out.(*os.File)
	if inputOK && outputOK && term.IsTerminal(input.Fd()) && term.IsTerminal(output.Fd()) && (name == "" || !allProjects) {
		var repositories []tui.Repository
		for _, project := range projects {
			repositories = append(repositories, tui.Repository{Name: project.Name, Branch: project.Branch})
			if allProjects {
				selected[project.Name] = true
			}
		}
		form, accepted, err := tui.RunWorkspaceForm(app.In, app.Out, tui.WorkspaceForm{
			Name: name, Agent: agent, Prompt: prompt, Repositories: repositories,
			Selected: selected, ValidateName: validateName,
		})
		if err != nil {
			return err
		}
		if !accepted {
			app.note("Cancelled.")
			return nil
		}
		name, agent, prompt, selected = form.Name, form.Agent, form.Prompt, form.Selected
		formUsed = true
	}
	interactive := name == ""
	if interactive {
		app.In = bufio.NewReader(app.In)
		var err error
		name, err = app.promptLine("Sandbox name: ")
		if err != nil {
			return err
		}
	}
	if err := validateName(name); err != nil {
		return err
	}
	if interactive {
		if !agentSet {
			value, err := app.promptValue(fmt.Sprintf("Agent (claude/codex) [%s]: ", agent))
			if err != nil {
				return err
			}
			if value != "" {
				agent = value
			}
		}
		if !promptSet {
			var err error
			prompt, err = app.promptValue("Task prompt (optional): ")
			if err != nil {
				return err
			}
		}
	}
	if agent != "claude" && agent != "codex" {
		return fmt.Errorf("agent must be 'claude' or 'codex', got: %s", agent)
	}
	for _, dependency := range []string{"git", "tmux"} {
		if err := app.checkCommand(dependency); err != nil {
			return err
		}
	}
	if err := app.needNative(); err != nil {
		return err
	}
	root, err := app.canonicalWorkspaceRoot(true)
	if err != nil {
		return err
	}
	workspace := filepath.Join(root, name)
	if _, err := os.Lstat(workspace); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("workspace already exists: %s", workspace)
	}
	background := context.Background()
	if app.tmuxSessionExists(background, name) {
		return fmt.Errorf("tmux session already exists on the '%s' server: %s", app.config.TmuxSocket, name)
	}
	exists, err := app.nativeSandboxNamedExists(background, name, true)
	if err != nil {
		return fmt.Errorf("cannot verify whether Docker Sandbox '%s' exists", name)
	}
	if exists {
		return fmt.Errorf("Docker Sandbox already exists: %s", name)
	}

	if formUsed {
		// The form already selected the repositories.
	} else if allProjects {
		for _, project := range projects {
			selected[project.Name] = true
		}
	} else {
		names := make([]string, 0, len(projects))
		for _, project := range projects {
			names = append(names, project.Name)
		}
		selection, ok, err := app.chooseWithFZF(background, "Projects > ", names, true)
		if err != nil {
			return err
		}
		if !ok {
			app.note("Cancelled.")
			return nil
		}
		for _, project := range selection {
			selected[project] = true
		}
	}
	selectionCount := 0
	for _, chosen := range selected {
		if chosen {
			selectionCount++
		}
	}
	if selectionCount == 0 {
		return errors.New("at least one project must be selected")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChannel)
	var signalStatus atomic.Int32
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case received := <-signalChannel:
			switch received {
			case syscall.SIGHUP:
				signalStatus.Store(129)
			case syscall.SIGINT:
				signalStatus.Store(130)
			case syscall.SIGTERM:
				signalStatus.Store(143)
			}
			cancel()
		case <-done:
		}
	}()

	workspaceCreated := false
	sandboxAttempted := false
	tmuxAttempted := false
	complete := false
	defer func() {
		if !complete {
			app.cleanupFailedNew(name, workspace, workspaceCreated, sandboxAttempted, tmuxAttempted)
		}
		if status := int(signalStatus.Load()); status != 0 {
			returnErr = statusError{code: status}
		}
	}()

	app.note("Creating %s...", name)
	if err := os.Mkdir(workspace, 0o777); err != nil {
		return fmt.Errorf("cannot create workspace: %s", workspace)
	}
	workspaceCreated = true
	if err := writeMarker(workspace, name); err != nil {
		return fmt.Errorf("cannot mark workspace: %s", workspace)
	}
	app.note("- workspace created: %s", workspace)

	var cloneSteps []tui.CreationStep
	for _, project := range projects {
		if selected[project.Name] {
			cloneSteps = append(cloneSteps, tui.CreationStep{Name: project.Name, Run: func() error {
				return app.cloneProjectFromCache(ctx, project, filepath.Join(workspace, project.Name))
			}})
		}
	}
	if output, ok := app.Out.(*os.File); ok && term.IsTerminal(output.Fd()) {
		var logs creationLog
		previousRunner, previousErr := app.runner, app.Err
		app.runner.Out, app.runner.Err, app.Err = &logs, &logs, &logs
		err := tui.RunCreation(app.Out, name, cloneSteps)
		app.runner, app.Err = previousRunner, previousErr
		if err != nil {
			fmt.Fprint(app.Err, logs.String())
			return err
		}
	} else {
		for _, step := range cloneSteps {
			if err := step.Run(); err != nil {
				return err
			}
		}
	}

	for _, project := range projects {
		if selected[project.Name] {
			if err := app.bootstrapProject(ctx, project, filepath.Join(workspace, project.Name)); err != nil {
				return err
			}
		}
	}

	sandboxAttempted = true
	app.note("- creating Docker Sandbox %s (%s)", name, agent)
	if err := app.runner.Run(ctx, app.native, []string{"create", "--name", name, agent, workspace}, process.Options{}); err != nil {
		return fmt.Errorf("Docker Sandbox creation failed: %s", name)
	}
	exists, err = app.nativeSandboxNamedExists(ctx, name, true)
	if err != nil || !exists {
		return fmt.Errorf("Docker Sandbox creation did not produce '%s'", name)
	}

	if agent == "claude" && app.config.ClaudeTheme != "" {
		if err := app.initializeClaudeTheme(ctx, name); err != nil {
			return err
		}
	}

	tmuxAttempted = true
	app.note("- creating isolated tmux session %s", name)
	if err := app.createAgentTmuxSession(ctx, name, workspace, name, prompt); err != nil {
		return fmt.Errorf("agent tmux session creation failed: %s", name)
	}
	complete = true
	app.note("Created sandbox workspace '%s'.", name)
	if !noOpen {
		return app.openSession(context.Background(), name)
	}
	return nil
}

// Git may write stdout and stderr concurrently.
type creationLog struct {
	sync.Mutex
	bytes.Buffer
}

func (log *creationLog) Write(data []byte) (int, error) {
	log.Lock()
	defer log.Unlock()
	return log.Buffer.Write(data)
}
