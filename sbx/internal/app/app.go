package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"sbx/internal/config"
	"sbx/internal/process"
)

type App struct {
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
	Program string
	Self    string
	BaseDir string

	runner       process.Runner
	config       config.Config
	configLoaded bool
	native       string
}

type statusError struct {
	code int
	err  error
}

func (err statusError) Error() string {
	if err.err == nil {
		return ""
	}
	return err.err.Error()
}

func New(executable string, in io.Reader, out, errorOutput io.Writer) (*App, error) {
	self, err := process.CanonicalExecutable(executable)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve executable path: %s", executable)
	}
	baseDir := filepath.Dir(filepath.Dir(self))
	return &App{
		In:      in,
		Out:     out,
		Err:     errorOutput,
		Program: filepath.Base(executable),
		Self:    self,
		BaseDir: baseDir,
		runner:  process.Runner{In: in, Out: out, Err: errorOutput},
	}, nil
}

func (app *App) Run(args []string) int {
	command := "help"
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	var err error
	switch command {
	case "-h", "--help", "help":
		app.usage()
		return 0
	case "new", "open", "add", "project", "split", "delete", "list", "doctor", "tui", "tmux-reload":
		if err = app.loadConfig(); err == nil {
			switch command {
			case "new":
				err = app.cmdNew(args)
			case "open":
				err = app.cmdOpen(args)
			case "add":
				err = app.cmdAdd(args)
			case "project":
				err = app.cmdProject(args)
			case "split":
				err = app.cmdSplit(args)
			case "delete":
				err = app.cmdDelete(args)
			case "list":
				err = app.cmdList(args)
			case "doctor":
				err = app.cmdDoctor(args)
			case "tui":
				err = app.cmdTUI(args)
			case "tmux-reload":
				err = app.cmdTmuxReload(args)
			}
		}
	case "__popup":
		err = app.cmdPopup(args)
	case "__sidecar-key":
		if err = app.loadConfig(); err == nil {
			err = app.cmdSidecarKey(args)
		}
	case "native":
		if len(args) == 0 {
			err = fmt.Errorf("Usage: %s native COMMAND [ARGUMENTS]", app.Program)
			break
		}
		if err = app.loadConfig(); err == nil {
			if err = app.needNative(); err == nil {
				err = process.Replace(app.native, args)
			}
		}
	default:
		if err = app.loadConfig(); err == nil {
			if err = app.needNative(); err == nil {
				err = process.Replace(app.native, append([]string{command}, args...))
			}
		}
	}

	if err == nil {
		return 0
	}
	var requestedStatus statusError
	if errors.As(err, &requestedStatus) {
		if requestedStatus.err != nil && requestedStatus.err.Error() != "" {
			fmt.Fprintf(app.Err, "error: %s\n", requestedStatus.err)
		}
		return requestedStatus.code
	}
	fmt.Fprintf(app.Err, "error: %s\n", err)
	return 1
}

func (app *App) loadConfig() error {
	if app.configLoaded {
		return nil
	}
	loaded, err := config.Load(app.BaseDir)
	if err != nil {
		return err
	}
	app.config = loaded
	app.configLoaded = true
	return nil
}

func (app *App) needNative() error {
	if app.native != "" {
		return nil
	}
	requested := os.Getenv("SBX_DOCKER_BIN")
	if requested == "" {
		requested = app.config.DockerSBXBinary
	}
	if requested != "" {
		var resolved string
		var err error
		if strings.ContainsRune(requested, filepath.Separator) {
			if !process.IsExecutable(requested) {
				return fmt.Errorf("Docker Sandboxes CLI is not executable: %s", requested)
			}
			resolved, err = process.CanonicalExecutable(requested)
			if err != nil {
				return fmt.Errorf("cannot resolve Docker Sandboxes CLI: %s", requested)
			}
		} else {
			resolved, err = exec.LookPath(requested)
			if err != nil {
				return fmt.Errorf("Docker Sandboxes CLI not found: %s", requested)
			}
			resolved, err = process.CanonicalExecutable(resolved)
			if err != nil {
				return fmt.Errorf("cannot resolve Docker Sandboxes CLI: %s", requested)
			}
		}
		if resolved == app.Self {
			return errors.New("docker_sbx_bin resolves to this workflow wrapper; point it at Docker's native sbx executable")
		}
		app.native = resolved
		return nil
	}

	resolved, err := process.FindNextOnPath("sbx", app.Self)
	if err != nil {
		return errors.New("Docker Sandboxes CLI not found (install docker/tap/sbx, or set docker_sbx_bin/SBX_DOCKER_BIN)")
	}
	app.native = resolved
	return nil
}

func (app *App) checkCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("required command not found: %s", name)
	}
	return nil
}

func (app *App) note(format string, values ...any) {
	fmt.Fprintf(app.Err, format+"\n", values...)
}

func (app *App) usage() {
	fmt.Fprintf(app.Out, `Usage: %s COMMAND [ARGUMENTS]

Commands:
  new [--all] [--no-open] [--agent claude|codex] [--prompt text] [name]
                            Select projects and create a sandbox workspace
  open [name]             Open or switch to a managed workspace
  add [project]          Clone a configured project into the current workspace
  project [name]          Open/switch a host project session
  split [project|agent]  Open a workspace session beside the current pane
  delete [--yes] [name]  Delete all resources for a workspace
  list                    List managed workspaces
  tui [--current] [--split]
                          Browse, optionally selecting a side-by-side session
  tmux-reload             Reload the user config and reapply sbx bindings
  doctor                  Validate configuration and dependencies
  native COMMAND [...]    Run Docker's native sbx explicitly
  help                    Show this help

This wrapper intentionally shares the name "sbx" with Docker Sandboxes. Set
docker_sbx_bin in config.yaml (or SBX_DOCKER_BIN), or put this wrapper before
Docker's native sbx on PATH; it will discover the next executable. Commands
other than the workflow commands above are forwarded to Docker's native sbx.
`, app.Program)
}

func (app *App) run(ctx context.Context, name string, args ...string) error {
	return app.runner.Run(ctx, name, args, process.Options{})
}
