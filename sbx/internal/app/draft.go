package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// newDraft keeps the form input of a failed `sbx new`, so the next attempt can
// restore it instead of the user retyping the task prompt.
type newDraft struct {
	Name     string   `json:"name"`
	Agent    string   `json:"agent"`
	Prompt   string   `json:"prompt"`
	Selected []string `json:"selected"`
}

func (app *App) newDraftPath() string {
	return filepath.Join(app.config.ConfigDir, "new-draft.json")
}

func (app *App) loadNewDraft() (newDraft, bool) {
	var draft newDraft
	if app.config.ConfigDir == "" {
		return draft, false
	}
	contents, err := os.ReadFile(app.newDraftPath())
	if err != nil || json.Unmarshal(contents, &draft) != nil {
		return newDraft{}, false
	}
	return draft, true
}

func (app *App) saveNewDraft(name, agent, prompt string, selected map[string]bool) error {
	if app.config.ConfigDir == "" {
		return errors.New("no configuration directory")
	}
	draft := newDraft{Name: name, Agent: agent, Prompt: prompt}
	for project, chosen := range selected {
		if chosen {
			draft.Selected = append(draft.Selected, project)
		}
	}
	sort.Strings(draft.Selected)
	contents, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return err
	}
	temp := app.newDraftPath() + ".next"
	if err := os.WriteFile(temp, append(contents, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(temp, app.newDraftPath())
}

func (app *App) clearNewDraft() {
	if app.config.ConfigDir != "" {
		_ = os.Remove(app.newDraftPath())
	}
}
