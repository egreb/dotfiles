package app

import (
	"context"
	"fmt"
	"os"
	"sbx/internal/model"
	"sbx/internal/process"
	"sbx/internal/tui"
)

func previousWorkspace(inventory model.Inventory, deleted string) string {
	var best *model.Workspace
	for i := range inventory.Workspaces {
		w := &inventory.Workspaces[i]
		if w.Name == deleted {
			continue
		}
		if best == nil || w.LastUsed > best.LastUsed || (w.LastUsed == best.LastUsed && w.Name < best.Name) {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	return best.Name
}

// Move the popup's client off sessions that deletion will close. If this is
// the last workspace, switch to an ordinary host shell session, whose name is
// returned so callers can clean it up.
func (app *App) leaveDeletedWorkspace(ctx context.Context, name, previous string) (string, error) {
	current, err := app.currentManagedContext(ctx)
	if err != nil || current.WorkspaceName != name || os.Getenv("SBX_TMUX_CLIENT") == "" {
		return "", nil
	}
	if previous != "" {
		return "", app.openSession(ctx, previous)
	}
	holding := fmt.Sprintf("sbx-picker-%d", os.Getpid())
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if err = app.tmuxRun(ctx, "new-session", "-d", "-s", holding, "-c", home, "exec "+process.ShellQuote(app.config.HostShell)+" -l"); err != nil {
		return "", err
	}
	if err = app.switchCurrentClient(ctx, holding); err != nil {
		_ = app.tmuxRunQuiet(ctx, "kill-session", "-t", "="+holding)
		return "", err
	}
	return holding, nil
}

// deleteCommand runs in its own tmux window, so failures stay on screen until
// acknowledged and successes are visible briefly before the window closes.
func deleteCommand(self, name string) string {
	quoted := process.ShellQuote(name)
	return "printf 'Deleting workspace %s…\\n\\n' " + quoted +
		"; " + process.ShellQuote(self) + " __popup delete --yes " + quoted +
		" && { printf '\\nWorkspace %s deleted.\\n' " + quoted + "; sleep 2; }; true"
}

// StartDelete removes a workspace from a tmux window beside the picker, so the
// picker popup keeps running. It reports false when there is no popup client
// to keep, and the picker then falls back to deleting in the foreground.
func (app *App) StartDelete(ctx context.Context, name string) (bool, error) {
	client := os.Getenv("SBX_TMUX_CLIENT")
	if client == "" {
		return false, nil
	}
	if err := validateName(name); err != nil {
		return true, err
	}
	inventory, err := app.Inventory(ctx)
	if err != nil {
		return true, err
	}
	if _, err = app.leaveDeletedWorkspace(ctx, name, previousWorkspace(inventory, name)); err != nil {
		return true, err
	}
	out, err := app.tmuxOutput(ctx, false, "display-message", "-p", "-c", client, "#{session_name}")
	if err != nil {
		return true, err
	}
	session := trimCommandOutput(out)
	return true, app.tmuxRun(ctx, "new-window", "-t", "="+session+":", "-n", "rm "+name, deleteCommand(app.Self, name))
}

// Called only after an explicit confirmation in the picker. The command-line
// delete flow retains its separate exact-name confirmation.
func (app *App) deleteFromPicker(name string) error {
	ctx := context.Background()
	inventory, err := app.Inventory(ctx)
	if err != nil {
		return err
	}
	previous := previousWorkspace(inventory, name)
	holding, err := app.leaveDeletedWorkspace(ctx, name, previous)
	if err != nil {
		return err
	}
	if holding != "" {
		defer func() {
			// A canceled creation or no remaining workspace leaves a usable shell.
			for _, s := range app.listTmuxSessions(ctx) {
				if s.Name == holding && !s.Attached {
					_ = app.tmuxRunQuiet(ctx, "kill-session", "-t", "="+holding)
				}
			}
		}()
	}
	if err = app.cmdDelete([]string{"--yes", name}); err != nil {
		return err
	}
	create, err := tui.Confirm(app.In, app.Out, "Workspace "+name+" deleted.\nCreate a new workspace?")
	if err != nil {
		return err
	}
	if create {
		return app.cmdNew(nil)
	}
	if previous != "" {
		return app.openSession(ctx, previous)
	}
	app.note("No workspaces remain. Use sbx new to create one.")
	return nil
}
