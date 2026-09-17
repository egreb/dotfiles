package tui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func formKey(m workspaceForm, key tea.KeyPressMsg) workspaceForm {
	next, _ := m.Update(key)
	return next.(workspaceForm)
}
func TestWorkspaceFormRevisitAndSelect(t *testing.T) {
	m := newWorkspaceForm(WorkspaceForm{Name: "test", Agent: "claude", Repositories: []Repository{{Name: "api"}, {Name: "web"}}})
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyRight})
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = formKey(m, keyPress("Fix checkout"))
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = formKey(m, keyPress(" "))
	if !m.data.Selected["api"] {
		t.Fatal("Space did not select repository")
	}
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.data.Selected["web"] {
		t.Fatal("Tab toggled a repository")
	}
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	if m.focus != 3 || m.data.Prompt != "Fix checkout" || m.data.Agent != "codex" {
		t.Fatalf("lost edits: %#v", m)
	}
	m = formKey(m, keyPress("web"))
	m = formKey(m, keyPress(" "))
	if !m.data.Selected["api"] || !m.data.Selected["web"] {
		t.Fatal("filter lost selection")
	}
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !m.accepted {
		t.Fatal("valid form did not submit")
	}
}
func TestWorkspaceFormValidationAndCancel(t *testing.T) {
	m := newWorkspaceForm(WorkspaceForm{Agent: "claude", ValidateName: func(s string) error {
		if s == "" {
			return errors.New("Name required")
		}
		return nil
	}, Repositories: []Repository{{Name: "api"}}})
	m.focus = 4
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.accepted || m.focus != 0 || m.err == "" {
		t.Fatal("invalid name submitted")
	}
	m.data.Name = "valid"
	m.focus = 4
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.accepted || m.focus != 3 {
		t.Fatal("empty repo selection submitted")
	}
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.accepted {
		t.Fatal("cancel submitted")
	}
}
func TestWorkspaceFormPasteAndEdit(t *testing.T) {
	m := newWorkspaceForm(WorkspaceForm{Agent: "claude"})
	m.focus = 2
	next, _ := m.Update(tea.PasteMsg{Content: "Fix\n日本語"})
	m = next.(workspaceForm)
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyLeft})
	m = formKey(m, tea.KeyPressMsg{Code: tea.KeyBackspace})
	if m.data.Prompt != "Fix\n日語" {
		t.Fatal(m.data.Prompt)
	}
}
func TestWorkspaceFormFitsPopup(t *testing.T) {
	m := newWorkspaceForm(WorkspaceForm{Agent: "claude", Repositories: []Repository{{Name: "api"}, {Name: "web"}, {Name: "cli"}, {Name: "docs"}, {Name: "other"}, {Name: "last"}}})
	m.width = 70
	m.height = 30
	view := ansi.Strip(m.View().Content)
	if strings.Count(view, "\n")+1 > m.height {
		t.Fatalf("form exceeds height: %d\n%s", strings.Count(view, "\n")+1, view)
	}
	if !strings.Contains(view, "selected to clone") {
		t.Fatal(view)
	}
}

func keyPress(s string) tea.KeyPressMsg { r := []rune(s); return tea.KeyPressMsg{Code: r[0], Text: s} }
