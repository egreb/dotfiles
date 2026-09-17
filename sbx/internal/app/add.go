package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"sbx/internal/config"
)

func (app *App) cmdAdd(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("Usage: %s add [project]", app.Program)
	}
	ctx := context.Background()
	current, err := app.currentManagedContext(ctx)
	if err != nil {
		return err
	}
	workspace, err := app.assertSafeWorkspace(current.Workspace, current.WorkspaceName)
	if err != nil {
		return err
	}
	projects, err := app.loadProjects()
	if err != nil {
		return err
	}
	name := ""
	if len(args) == 1 {
		name = args[0]
	}
	if name == "" {
		var choices []string
		for _, project := range projects {
			_, err := os.Lstat(filepath.Join(workspace, project.Name))
			if os.IsNotExist(err) {
				choices = append(choices, project.Name)
			} else if err != nil {
				return err
			}
		}
		if len(choices) == 0 {
			app.note("All configured projects are already in this workspace.")
			return nil
		}
		selection, ok, err := app.chooseWithFZF(ctx, "Add project > ", choices, false)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		name = selection[0]
	}
	var selected *config.Project
	for i := range projects {
		if projects[i].Name == name {
			selected = &projects[i]
			break
		}
	}
	if selected == nil {
		return fmt.Errorf("project is not configured: %s", name)
	}
	if err := validateName(name); err != nil {
		return err
	}
	if err := app.checkCommand("git"); err != nil {
		return err
	}
	destination := filepath.Join(workspace, name)
	// Reserve the destination atomically: never clone over files or symlinks.
	if err := os.Mkdir(destination, 0755); err != nil {
		return fmt.Errorf("cannot add project (path may already exist): %w", err)
	}
	if err := app.cloneProjectFromCache(ctx, *selected, destination); err != nil {
		// Remove only an empty reservation; retain any partial clone for inspection.
		_ = os.Remove(destination)
		return fmt.Errorf("%w; inspect %s before retrying", err, destination)
	}
	if err := app.bootstrapProject(ctx, *selected, destination); err != nil {
		return fmt.Errorf("project cloned at %s but setup failed: %w", destination, err)
	}
	app.note("Added %s to workspace %s.", name, current.WorkspaceName)
	return nil
}
