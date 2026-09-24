package tui

import (
	"sbx/internal/model"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestLongWorkspaceAndProjectRemainReadable(t *testing.T) {
	workspace := "presence-v2-remove-working-hours"
	project := "baloo-frontend-platform"
	for _, width := range []int{45, 80, 120, 220} {
		m := dashboard{width: width, height: 26, showPane: true, inventory: model.Inventory{Workspaces: []model.Workspace{{Name: workspace, Projects: []model.Project{{Name: project}}}}}}
		m.filter()
		m.cursor = 1
		content := ansi.Strip(m.View().Content)
		joined := strings.ReplaceAll(content, "\n", "")
		if !strings.Contains(joined, workspace) || !strings.Contains(joined, project) {
			t.Fatalf("%d: selected identity missing: %s", width, content)
		}
		for _, line := range strings.Split(content, "\n") {
			if ansi.StringWidth(line) > width {
				t.Fatalf("%d: overflow: %s", width, line)
			}
		}
		if len(strings.Split(content, "\n")) > m.height {
			t.Fatal("height overflow")
		}
	}
}
func TestDashboardFilterAndStalePreview(t *testing.T) {
	m := dashboard{width: 100, height: 25, inventory: model.Inventory{Workspaces: []model.Workspace{{Name: "alpha", AgentSession: &model.Session{Name: "alpha"}, Projects: []model.Project{{Name: "backend", Session: &model.Session{Name: "alpha--backend"}}, {Name: "frontend"}}}, {Name: "beta"}}}}
	m.filter()
	m.cursor = 1
	m.query = "ALP BK"
	m.filter()
	if len(m.rows) != 1 || m.selected().selection.Project != "backend" {
		t.Fatalf("%+v", m.rows)
	}
	updated, _ := m.Update(previewMsg{session: "beta", content: "wrong preview"})
	if updated.(dashboard).preview != "" {
		t.Fatal("stale preview replaced selected preview")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if updated.(dashboard).action.Selection.Project != "backend" {
		t.Fatal("wrong selected action")
	}
	m.options = Options{Workspace: "alpha", ExcludedSession: "alpha"}
	m.query = ""
	m.filter()
	if len(m.rows) != 2 || strings.Contains(m.rows[0].label, "alpha") {
		t.Fatalf("scoped rows: %+v", m.rows)
	}
}

func TestDeletePickerSelectsWholeWorkspaces(t *testing.T) {
	m := dashboard{options: Options{DeleteMode: true}, inventory: model.Inventory{Workspaces: []model.Workspace{{Name: "alpha", Projects: []model.Project{{Name: "backend"}, {Name: "frontend"}}}, {Name: "beta"}}}}
	m.filter()
	if len(m.rows) != 2 {
		t.Fatalf("expected one row per workspace: %+v", m.rows)
	}
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if updated.(dashboard).action.Kind != ActionNone {
		t.Fatal("escape requested deletion")
	}
	m.query = "beta"
	m.filter()
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	action := updated.(dashboard).action
	if action.Kind != ActionDelete || action.Selection.Workspace != "beta" || action.Selection.Project != "" {
		t.Fatalf("wrong deletion selection: %+v", action)
	}
}
