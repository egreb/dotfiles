package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sbx/internal/process"
	"strings"
)

type tmuxSession struct {
	Name           string
	Managed        bool
	WorkspaceName  string
	Workspace      string
	DockerName     string
	Role           string
	Project        string
	LastSession    string
	Attached       bool
	ActivePaneID   string
	CurrentCommand string
	PaneDead       bool
}

type managedContext struct {
	Session       string
	Pane          string
	Workspace     string
	DockerName    string
	WorkspaceName string
}

func (app *App) tmuxArgs(args ...string) []string {
	prefix := []string{"-L", app.config.TmuxSocket}
	if app.config.TmuxConfig != "" {
		prefix = append(prefix, "-f", app.config.TmuxConfig)
	}
	return append(prefix, args...)
}

func (app *App) tmuxRun(ctx context.Context, args ...string) error {
	return app.runner.Run(ctx, "tmux", app.tmuxArgs(args...), process.Options{})
}

func (app *App) tmuxRunQuiet(ctx context.Context, args ...string) error {
	return app.runner.Run(ctx, "tmux", app.tmuxArgs(args...), process.Options{Stdout: io.Discard, QuietErr: true})
}

func (app *App) tmuxOutput(ctx context.Context, quiet bool, args ...string) (string, error) {
	return app.runner.Output(ctx, "tmux", app.tmuxArgs(args...), process.Options{QuietErr: quiet})
}

func (app *App) tmuxSessionExists(ctx context.Context, wanted string) bool {
	output, _ := app.tmuxOutput(ctx, true, "list-sessions", "-F", "#{session_name}")
	for _, session := range outputLines(output) {
		if session == wanted {
			return true
		}
	}
	return false
}

func (app *App) setTmuxWorkspaceStatus(ctx context.Context, session, workspaceName string) error {
	return app.tmuxRun(ctx, "set-option", "-t", session, "status-left", fmt.Sprintf("#[bold] sbx:%s #[default]│ ", workspaceName))
}

func (app *App) ensureProjectSession(ctx context.Context, current managedContext, name string, remember bool) (string, error) {
	projectSession := current.WorkspaceName + "--" + name
	if app.tmuxSessionExists(ctx, projectSession) {
		checks := []struct {
			actual   string
			expected string
			message  string
		}{
			{app.tmuxOption(ctx, projectSession, "@sbx-managed"), "1", "refusing to use unmanaged tmux session"},
			{app.tmuxOption(ctx, projectSession, "@sbx-name"), current.WorkspaceName, "project session belongs to another workspace"},
			{app.tmuxOption(ctx, projectSession, "@sbx-role"), "project", "tmux session is not a project session"},
			{app.tmuxOption(ctx, projectSession, "@sbx-project"), name, "project session metadata does not match"},
			{app.tmuxOption(ctx, projectSession, "@sbx-workspace"), current.Workspace, "project session has the wrong workspace path"},
			{app.tmuxOption(ctx, projectSession, "@sbx-docker-name"), current.DockerName, "project session has the wrong Docker Sandbox"},
		}
		for _, check := range checks {
			if check.actual != check.expected {
				return "", fmt.Errorf("%s: %s", check.message, projectSession)
			}
		}
		if err := app.setTmuxWorkspaceStatus(ctx, projectSession, current.WorkspaceName); err != nil {
			return "", fmt.Errorf("cannot set tmux workspace status: %s", projectSession)
		}
		if remember {
			if err := app.setWorkspaceLastSession(ctx, current.WorkspaceName, projectSession); err != nil {
				return "", err
			}
		}
		return projectSession, nil
	}

	hostCommand := fmt.Sprintf("exec %s -l", process.ShellQuote(app.config.HostShell))
	if err := app.tmuxRun(ctx, "new-session", "-d", "-s", projectSession, "-n", "shell", "-c", filepath.Join(current.Workspace, name), hostCommand); err != nil {
		return "", fmt.Errorf("cannot create project tmux session: %s", projectSession)
	}
	configured := false
	defer func() {
		if !configured {
			_ = app.tmuxRunQuiet(context.Background(), "kill-session", "-t", projectSession)
		}
	}()
	settings := [][2]string{
		{"@sbx-managed", "1"},
		{"@sbx-name", current.WorkspaceName},
		{"@sbx-workspace", current.Workspace},
		{"@sbx-docker-name", current.DockerName},
		{"@sbx-role", "project"},
		{"@sbx-project", name},
	}
	for _, setting := range settings {
		if err := app.tmuxRun(ctx, "set-option", "-t", projectSession, setting[0], setting[1]); err != nil {
			return "", fmt.Errorf("cannot configure project tmux session: %s", projectSession)
		}
	}
	if err := app.tmuxRun(ctx, "set-window-option", "-t", projectSession, "automatic-rename", "on"); err != nil {
		return "", fmt.Errorf("cannot configure project tmux session: %s", projectSession)
	}
	if err := app.setTmuxWorkspaceStatus(ctx, projectSession, current.WorkspaceName); err != nil {
		return "", fmt.Errorf("cannot configure project tmux session: %s", projectSession)
	}
	if err := app.configureTmuxEnvironment(ctx); err != nil {
		return "", fmt.Errorf("cannot configure project tmux session: %s", projectSession)
	}
	if err := app.initializeProjectWindows(ctx, projectSession, filepath.Join(current.Workspace, name)); err != nil {
		return "", err
	}
	configured = true
	if remember {
		if err := app.setWorkspaceLastSession(ctx, current.WorkspaceName, projectSession); err != nil {
			return "", err
		}
	}
	return projectSession, nil
}

func (app *App) initializeProjectWindows(ctx context.Context, session, directory string) error {
	target := "=" + session + ":"
	if err := app.tmuxRun(ctx, "set-option", "-t", session, "base-index", "1"); err != nil {
		return err
	}
	index, err := app.tmuxOutput(ctx, false, "display-message", "-p", "-t", target, "#{window_index}")
	if err != nil {
		return err
	}
	if trimCommandOutput(index) != "1" {
		if err := app.tmuxRun(ctx, "move-window", "-s", target, "-t", target+"1"); err != nil {
			return err
		}
	}
	hostCommand := fmt.Sprintf("exec %s -l", process.ShellQuote(app.config.HostShell))
	for i, tool := range [][2]string{{"neovim", "env NVIM_APPNAME=nvim-v2 nvim"}, {"lazygit", "lazygit"}, {"hunk", "hunk diff"}} {
		window := fmt.Sprintf("%s%d", target, i+1)
		if i > 0 {
			if err := app.tmuxRun(ctx, "new-window", "-d", "-t", window, "-n", tool[0], "-c", directory, hostCommand); err != nil {
				return err
			}
		}
		if err := app.tmuxRun(ctx, "set-window-option", "-t", window, "automatic-rename", "off"); err != nil {
			return err
		}
		if err := app.tmuxRun(ctx, "rename-window", "-t", window, tool[0]); err != nil {
			return err
		}
		if err := app.tmuxRun(ctx, "send-keys", "-t", window, "-l", tool[1]); err != nil {
			return err
		}
		if err := app.tmuxRun(ctx, "send-keys", "-t", window, "Enter"); err != nil {
			return err
		}
	}
	return app.tmuxRun(ctx, "select-window", "-t", target+"1")
}

func (app *App) attachOrSwitch(ctx context.Context, target string) error {
	// Return popup selections to the invoking client, even without TMUX in its environment.
	if os.Getenv("SBX_TMUX_CLIENT") != "" {
		return app.switchCurrentClient(ctx, target)
	}
	socketOutput, _ := app.tmuxOutput(ctx, true, "display-message", "-p", "#{socket_path}")
	if tmuxEnvironment := os.Getenv("TMUX"); tmuxEnvironment != "" && firstTmuxField(tmuxEnvironment) == trimCommandOutput(socketOutput) {
		return app.switchCurrentClient(ctx, target)
	}
	return app.runner.Run(ctx, "tmux", app.tmuxArgs("attach-session", "-t", target), process.Options{
		Env: map[string]string{"TMUX": "", "TMUX_PANE": ""},
	})
}

func outputLines(output string) []string {
	output = strings.TrimRight(output, "\n")
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}

func trimCommandOutput(output string) string {
	return strings.TrimRight(output, "\n")
}

func firstTmuxField(value string) string {
	if index := strings.IndexByte(value, ','); index >= 0 {
		return value[:index]
	}
	return value
}
