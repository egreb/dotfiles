package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sbx/internal/model"
)

func (app *App) Inventory(ctx context.Context) (model.Inventory, error) {
	projects, err := app.loadProjects()
	if err != nil {
		return model.Inventory{}, err
	}
	candidates, err := app.workspaceCandidates()
	if err != nil {
		return model.Inventory{}, err
	}
	sessions := app.listTmuxSessions(ctx)
	sessionsByName := make(map[string]tmuxSession, len(sessions))
	for _, session := range sessions {
		if session.Managed {
			sessionsByName[session.Name] = session
		}
	}

	// Inventory is deliberately filesystem- and tmux-only. Calling the native
	// Docker CLI can initiate an OS security prompt, so it must happen only for
	// an explicit create, recover, or delete action.
	inventory := model.Inventory{Workspaces: make([]model.Workspace, 0, len(candidates))}
	for _, name := range candidates {
		workspacePath := filepath.Join(app.config.WorkspaceRoot, name)
		workspace := model.Workspace{
			Name:     name,
			Path:     workspacePath,
			Projects: make([]model.Project, 0),
		}
		if session, ok := sessionsByName[name]; ok && session.Role != "project" {
			workspace.AgentSession = modelSession(session)
		}
		for _, projectConfig := range projects {
			projectPath := filepath.Join(workspacePath, projectConfig.Name)
			if info, statErr := os.Stat(projectPath); statErr != nil || !info.IsDir() {
				continue
			}
			project := model.Project{Name: projectConfig.Name, Path: projectPath}
			if session, ok := sessionsByName[name+"--"+projectConfig.Name]; ok && session.Role == "project" && session.WorkspaceName == name {
				project.Session = modelSession(session)
			}
			workspace.Projects = append(workspace.Projects, project)
		}
		for _, session := range sessionsByName {
			if session.WorkspaceName == name && session.LastUsed > workspace.LastUsed {
				workspace.LastUsed = session.LastUsed
			}
		}
		inventory.Workspaces = append(inventory.Workspaces, workspace)
	}
	sort.SliceStable(inventory.Workspaces, func(left, right int) bool {
		return inventory.Workspaces[left].Name < inventory.Workspaces[right].Name
	})
	return inventory, nil
}

func modelSession(session tmuxSession) *model.Session {
	return &model.Session{
		Name:           session.Name,
		Role:           session.Role,
		Project:        session.Project,
		Attached:       session.Attached,
		ActivePaneID:   session.ActivePaneID,
		CurrentCommand: session.CurrentCommand,
		PaneDead:       session.PaneDead,
	}
}

func (app *App) Preview(ctx context.Context, session string, lines int) (string, error) {
	if session == "" {
		return "No live tmux session selected.", nil
	}
	if lines < 1 {
		lines = 1
	}
	output, err := app.tmuxOutput(ctx, true, "capture-pane", "-p", "-t", session+":", "-S", fmt.Sprintf("-%d", lines))
	if err != nil {
		return "", fmt.Errorf("session preview unavailable")
	}
	return sanitizePreview(strings.TrimRight(output, "\n")), nil
}

func sanitizePreview(value string) string {
	var result strings.Builder
	runes := []rune(value)
	for index := 0; index < len(runes); index++ {
		character := runes[index]
		if character == 0x1b {
			if index+1 < len(runes) && runes[index+1] == '[' {
				index += 2
				for index < len(runes) && (runes[index] < '@' || runes[index] > '~') {
					index++
				}
			} else if index+1 < len(runes) {
				index++
			}
			continue
		}
		switch {
		case character == '\n' || character == '\t':
			result.WriteRune(character)
		case character >= 0x20 && character != 0x7f:
			result.WriteRune(character)
		}
	}
	return result.String()
}

func (app *App) activateSelection(ctx context.Context, selection model.Selection) error {
	if selection.Kind == model.SelectWorkspace {
		return app.openSession(ctx, selection.Workspace)
	}
	if selection.Kind == model.SelectAgent {
		if !app.tmuxSessionManaged(ctx, selection.Workspace) {
			if err := app.restoreAgentTmuxSession(ctx, selection.Workspace); err != nil {
				return err
			}
		}
		if err := app.setTmuxWorkspaceStatus(ctx, selection.Workspace, selection.Workspace); err != nil {
			return fmt.Errorf("cannot set tmux workspace status: %s", selection.Workspace)
		}
		if err := app.setWorkspaceLastSession(ctx, selection.Workspace, selection.Workspace); err != nil {
			return err
		}
		return app.attachOrSwitch(ctx, selection.Workspace)
	}
	if selection.Kind != model.SelectProject {
		return nil
	}
	if !app.tmuxSessionManaged(ctx, selection.Workspace) {
		if err := app.restoreAgentTmuxSession(ctx, selection.Workspace); err != nil {
			return err
		}
	}
	workspace := app.tmuxOption(ctx, selection.Workspace, "@sbx-workspace")
	dockerName := app.tmuxOption(ctx, selection.Workspace, "@sbx-docker-name")
	if workspace == "" || dockerName == "" {
		return fmt.Errorf("agent session metadata is incomplete: %s", selection.Workspace)
	}
	projects, err := app.loadProjects()
	if err != nil {
		return err
	}
	if !app.projectConfigured(selection.Project, projects) {
		return fmt.Errorf("project is not configured: %s", selection.Project)
	}
	if info, statErr := os.Stat(filepath.Join(workspace, selection.Project)); statErr != nil || !info.IsDir() {
		return fmt.Errorf("project is not cloned in this workspace: %s", selection.Project)
	}
	current := managedContext{WorkspaceName: selection.Workspace, Workspace: workspace, DockerName: dockerName}
	target, err := app.ensureProjectSession(ctx, current, selection.Project, true)
	if err != nil {
		return err
	}
	if err := app.setTmuxWorkspaceStatus(ctx, target, selection.Workspace); err != nil {
		return err
	}
	return app.attachOrSwitch(ctx, target)
}
