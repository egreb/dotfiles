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

// Called only after an explicit confirmation in the picker. The command-line
// delete flow retains its separate exact-name confirmation.
func (app *App) deleteFromPicker(name string) error {
	ctx := context.Background()
	inventory, err := app.Inventory(ctx)
	if err != nil {
		return err
	}
	previous := previousWorkspace(inventory, name)
	// Move the popup's client off sessions that deletion will close. If this is
	// the last workspace, keep an ordinary host shell so the follow-up survives.
	current, currentErr := app.currentManagedContext(ctx)
	if currentErr == nil && current.WorkspaceName == name && os.Getenv("SBX_TMUX_CLIENT") != "" {
		if previous != "" {
			if err = app.openSession(ctx, previous); err != nil {
				return err
			}
		} else {
			holding := fmt.Sprintf("sbx-picker-%d", os.Getpid())
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			if err = app.tmuxRun(ctx, "new-session", "-d", "-s", holding, "-c", home, "exec "+process.ShellQuote(app.config.HostShell)+" -l"); err != nil {
				return err
			}
			if err = app.switchCurrentClient(ctx, holding); err != nil {
				_ = app.tmuxRunQuiet(ctx, "kill-session", "-t", "="+holding)
				return err
			}
			defer func() {
				// A canceled creation or no remaining workspace leaves a usable shell.
				for _, s := range app.listTmuxSessions(ctx) {
					if s.Name == holding && !s.Attached {
						_ = app.tmuxRunQuiet(ctx, "kill-session", "-t", "="+holding)
					}
				}
			}()
		}
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
