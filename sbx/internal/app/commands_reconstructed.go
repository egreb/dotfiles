package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/model"
	"sbx/internal/process"
	"strings"
)

func (app *App) nativeSandboxNamedExists(ctx context.Context, name string, quiet bool) (bool, error) {
	if e := app.needNative(); e != nil {
		return false, e
	}
	out, e := app.runner.Output(ctx, app.native, []string{"ls", "--format", "{{.Name}}"}, process.Options{QuietErr: quiet})
	if e != nil {
		return false, e
	}
	for _, l := range outputLines(out) {
		if strings.TrimSpace(l) == name {
			return true, nil
		}
	}
	return false, nil
}
func (app *App) cleanupFailedNew(name, path string, created, sandbox, tmux bool) {
	ctx := context.Background()
	if tmux {
		for _, s := range app.listTmuxSessions(ctx) {
			if s.Managed && s.WorkspaceName == name && s.Workspace == path {
				_ = app.tmuxRunQuiet(ctx, "kill-session", "-t", "="+s.Name)
			}
		}
	}
	if sandbox {
		if e := app.runner.Run(ctx, app.native, []string{"rm", name}, process.Options{QuietErr: true}); e != nil {
			app.note("Sandbox cleanup failed; keeping %s for inspection", path)
			return
		}
	}
	if created {
		if safe, e := app.assertSafeWorkspace(path, name); e == nil {
			if e = os.RemoveAll(safe); e != nil {
				app.note("Cleanup failed: %v", e)
			}
		}
	}
}
func (app *App) cmdOpen(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("Usage: sbx open [name]")
	}
	if len(args) == 0 {
		return app.cmdTUI(nil)
	}
	return app.openSession(context.Background(), args[0])
}
func (app *App) cmdProject(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("Usage: sbx project [name]")
	}
	if len(args) == 0 {
		return app.cmdTUI([]string{"--current"})
	}
	c, e := app.currentManagedContext(context.Background())
	if e != nil {
		return e
	}
	return app.activateSelection(context.Background(), model.Selection{Kind: model.SelectProject, Workspace: c.WorkspaceName, Project: args[0]})
}
func (app *App) cmdList(args []string) error {
	if len(args) > 1 || len(args) == 1 && args[0] != "--names" {
		return fmt.Errorf("Usage: sbx list [--names]")
	}
	names, e := app.workspaceCandidates()
	if e != nil {
		return e
	}
	for _, name := range names {
		fmt.Fprintln(app.Out, name)
	}
	return nil
}
func (app *App) cmdDoctor(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Usage: sbx doctor")
	}
	if _, e := app.loadProjects(); e != nil {
		return e
	}
	for _, bin := range []string{"git", "tmux", app.config.HostShell} {
		if e := app.checkCommand(bin); e != nil {
			return e
		}
	}
	if e := app.needNative(); e != nil {
		return e
	}
	fmt.Fprintln(app.Out, "Configuration and required executables OK")
	return nil
}
func (app *App) cmdDelete(args []string) error {
	yes := false
	name := ""
	for _, v := range args {
		if v == "--yes" {
			yes = true
		} else if name == "" && !strings.HasPrefix(v, "-") {
			name = v
		} else {
			return fmt.Errorf("Usage: sbx delete [--yes] [name]")
		}
	}
	if name == "" {
		names, e := app.workspaceCandidates()
		if e != nil {
			return e
		}
		rows, ok, e := app.chooseWithFZF(context.Background(), "Delete workspace > ", names, false)
		if e != nil || !ok {
			return e
		}
		name = rows[0]
	}
	if e := validateName(name); e != nil {
		return e
	}
	path, e := app.assertSafeWorkspace(filepath.Join(app.config.WorkspaceRoot, name), name)
	if e != nil {
		return e
	}
	if !yes {
		answer, e := app.promptValue("Delete workspace " + name + " and its sandbox? Type the workspace name: ")
		if e != nil {
			return e
		}
		if answer != name {
			app.note("Cancelled.")
			return nil
		}
	}
	if e = app.needNative(); e != nil {
		return e
	}
	ctx := context.Background()
	exists, e := app.nativeSandboxNamedExists(ctx, name, false)
	if e != nil {
		return e
	}
	if exists {
		if e = app.runner.Run(ctx, app.native, []string{"rm", name}, process.Options{}); e != nil {
			return e
		}
	}
	for _, s := range app.listTmuxSessions(ctx) {
		if s.Managed && s.WorkspaceName == name && s.Workspace == path {
			if e = app.tmuxRun(ctx, "kill-session", "-t", "="+s.Name); e != nil {
				return e
			}
		}
	}
	safe, e := app.assertSafeWorkspace(path, name)
	if e != nil {
		return e
	}
	return os.RemoveAll(safe)
}
