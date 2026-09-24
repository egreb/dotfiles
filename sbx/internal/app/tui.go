package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sbx/internal/process"
	"sbx/internal/tui"
)

func (app *App) cmdTUI(args []string) error {
	currentOnly := false
	splitMode := false
	for _, argument := range args {
		switch argument {
		case "--current":
			currentOnly = true
		case "--split":
			splitMode = true
			currentOnly = true
		default:
			return fmt.Errorf("Usage: %s tui [--current] [--split]", app.Program)
		}
	}
	if err := app.checkCommand("tmux"); err != nil {
		return err
	}
	ctx := context.Background()
	options := tui.Options{}
	current, currentErr := app.currentManagedContext(ctx)
	if currentErr == nil {
		options.SelectedSession = current.Session
		if currentOnly {
			options.Workspace = current.WorkspaceName
		}
	} else if currentOnly {
		return currentErr
	}
	if splitMode {
		options.ExcludedSession = current.Session
		options.Title = "sbx split"
	}
	action, err := tui.Run(app, app.In, app.Out, options)
	if err != nil {
		return err
	}
	switch action.Kind {
	case tui.ActionNone:
		return nil
	case tui.ActionActivate:
		if splitMode {
			return app.splitSelection(context.Background(), current, action.Selection)
		}
		return app.activateSelection(context.Background(), action.Selection)
	case tui.ActionNew:
		return app.cmdNew(nil)
	case tui.ActionDeleteConfirmed:
		return app.deleteFromPicker(action.Selection.Workspace)
	case tui.ActionDelete:
		return app.cmdDelete([]string{action.Selection.Workspace})
	default:
		return nil
	}
}

func (app *App) cmdPopup(args []string) error {
	if len(args) == 0 {
		return errors.New("popup requires a workflow command")
	}
	action := args[0]
	if action != "add" && action != "new" && action != "open" && action != "project" && action != "tui" && action != "delete" {
		return fmt.Errorf("unsupported popup command: %s", action)
	}
	err := app.runner.Run(context.Background(), app.Self, append([]string{action}, args[1:]...), process.Options{})
	status := process.ExitCode(err)
	if status != 0 {
		fmt.Fprintf(app.Err, "\nCommand failed (exit %d). Press Enter to close this popup.", status)
		if terminal, openErr := os.Open("/dev/tty"); openErr == nil {
			_, _ = bufio.NewReader(terminal).ReadString('\n')
			_ = terminal.Close()
		}
		fmt.Fprintln(app.Err)
	}
	if status != 0 {
		return statusError{code: status}
	}
	return nil
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func lookPath(name string) (string, error) {
	return exec.LookPath(name)
}
