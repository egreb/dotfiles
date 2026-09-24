package tui

import (
	tea "charm.land/bubbletea/v2"
	"sbx/internal/model"
	"strings"
	"testing"
)

func treeFixture() dashboard {
	m := dashboard{width: 100, height: 30, inventory: model.Inventory{Workspaces: []model.Workspace{
		{Name: "alpha", LastUsed: 10, Projects: []model.Project{{Name: "api"}, {Name: "web"}}},
		{Name: "beta", LastUsed: 20, Projects: []model.Project{{Name: "worker"}}},
		{Name: "unknown"},
	}}}
	m.filter()
	return m
}
func press(m dashboard, key rune) dashboard {
	next, _ := m.Update(tea.KeyPressMsg{Code: key, Text: string(key)})
	return next.(dashboard)
}
func TestTreeNavigationStaysWithinWorkspaceAndRemembersSession(t *testing.T) {
	m := treeFixture()
	m = press(m, 'j')
	m = press(m, 'j')
	m = press(m, 'j')
	if m.selected().selection.Project != "web" {
		t.Fatal("j crossed workspace boundary")
	}
	m = press(m, 'l')
	if m.selected().selection.Workspace != "beta" {
		t.Fatal("l did not jump workspace")
	}
	m = press(m, 'k')
	if m.selected().selection.Workspace != "beta" {
		t.Fatal("k crossed workspace boundary")
	}
	m = press(m, 'h')
	if m.selected().selection.Project != "web" {
		t.Fatal("workspace selection not restored")
	}
	m = press(m, 'h')
	if m.selected().selection.Workspace != "alpha" {
		t.Fatal("h passed first workspace")
	}
}
func TestTreeSortingFilteringAndRefresh(t *testing.T) {
	m := treeFixture()
	m = press(m, 'j')
	m = press(m, 's')
	if m.rows[0].selection.Workspace != "beta" || m.selected().selection.Project != "api" {
		t.Fatal("recency sort or selection preservation failed")
	}
	if m.rows[len(m.rows)-1].selection.Workspace != "unknown" {
		t.Fatal("unknown must sort last")
	}
	m = press(m, '/')
	m = press(m, 'w')
	m = press(m, 'k')
	m = press(m, 'r')
	if m.query != "wkr" || m.selected().selection.Project != "worker" {
		t.Fatal("navigation keys should type in search")
	}
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = next.(dashboard)
	if m.searching || m.action.Kind != ActionNone {
		t.Fatal("enter should finish search without opening")
	}
	next, _ = m.Update(inventoryMsg{inventory: m.inventory})
	m = next.(dashboard)
	if m.selected().selection.Project != "worker" {
		t.Fatal("refresh lost selection")
	}
	next, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if next.(dashboard).action.Selection.Project != "worker" {
		t.Fatal("wrong activation")
	}
}
func TestTreeRendersGroupsAndFolds(t *testing.T) {
	m := treeFixture()
	lines, _ := m.treeLines(100)
	content := strings.Join(lines, "\n")
	if strings.Count(content, "alpha") != 1 || !strings.Contains(content, "api") || strings.Contains(content, "worker") {
		t.Fatal(content)
	}
	m = press(m, ' ')
	lines, _ = m.treeLines(100)
	if strings.Contains(strings.Join(lines, "\n"), "api") {
		t.Fatal("space did not fold")
	}
	m = press(m, 'j')
	lines, _ = m.treeLines(100)
	if !strings.Contains(strings.Join(lines, "\n"), "api") {
		t.Fatal("j did not expand")
	}
}

func TestPickerDeleteRequiresExplicitConfirmation(t *testing.T) {
	for _, cancel := range []string{"n", "N", "esc", "enter", "ctrl+c"} {
		m := treeFixture()
		m = press(m, 'j')
		m = press(m, 'D')
		if m.pendingDelete != "alpha" || m.action.Kind != ActionNone {
			t.Fatal("D must confirm selected project's workspace")
		}
		if !strings.Contains(m.View().Content, "Uncommitted changes will be lost") {
			t.Fatal("missing destructive consequences")
		}
		var key tea.KeyPressMsg
		switch cancel {
		case "esc":
			key = tea.KeyPressMsg{Code: tea.KeyEscape}
		case "enter":
			key = tea.KeyPressMsg{Code: tea.KeyEnter}
		case "ctrl+c":
			key = tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
		default:
			key = tea.KeyPressMsg{Code: rune(cancel[0]), Text: cancel}
		}
		next, _ := m.Update(key)
		m = next.(dashboard)
		if m.pendingDelete != "" || m.action.Kind != ActionNone || m.selected().selection.Project != "api" {
			t.Fatalf("cancel %s failed", cancel)
		}
	}
	m := treeFixture()
	m = press(m, 'D')
	m = press(m, 'l')
	if m.selected().selection.Workspace != "alpha" {
		t.Fatal("confirmation allowed navigation")
	}
	m = press(m, 'y')
	if m.action.Kind != ActionDeleteConfirmed || m.action.Selection.Workspace != "alpha" || m.action.Selection.Kind != model.SelectWorkspace {
		t.Fatal("wrong deletion")
	}
	m = treeFixture()
	m = press(m, '/')
	m = press(m, 'D')
	if m.pendingDelete != "" || m.query != "D" {
		t.Fatal("D in filter must type")
	}
}
func TestCreateConfirmationDefaultsToNo(t *testing.T) {
	m := confirmation{question: "Create a new workspace?", width: 80}
	for _, key := range []tea.KeyPressMsg{{Code: tea.KeyEnter}, {Code: tea.KeyEscape}, {Code: 'n', Text: "n"}} {
		next, cmd := m.Update(key)
		if next.(confirmation).accepted || cmd == nil {
			t.Fatal("confirmation must default to no")
		}
	}
	next, cmd := m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if !next.(confirmation).accepted || cmd == nil {
		t.Fatal("yes not accepted")
	}
}
