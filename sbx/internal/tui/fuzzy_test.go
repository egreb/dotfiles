package tui

import (
	"sbx/internal/model"
	"testing"
)

func TestFuzzyScorePrefersBoundariesAndRuns(t *testing.T) {
	for _, c := range []struct{ query, better, worse string }{
		{"api", "api", "aptitude"},
		{"fe", "frontend", "platform-service"},
		{"bf", "baloo-frontend", "backoff"},
		{"web", "web-client", "hotwebsocket"},
		{"svc", "baloo-svc", "basic-volume-cache"},
	} {
		b, okB := fuzzyScore(c.better, c.query)
		w, okW := fuzzyScore(c.worse, c.query)
		if !okB || !okW || b <= w {
			t.Errorf("%q: %q=%d should beat %q=%d", c.query, c.better, b, c.worse, w)
		}
	}
	if _, ok := fuzzyScore("backend", "bk x"); ok {
		t.Error("every term must match")
	}
}

func TestScopedFilterRanksProjectsAndIgnoresWorkspaceName(t *testing.T) {
	m := dashboard{options: Options{Workspace: "presence-v2"}, inventory: model.Inventory{Workspaces: []model.Workspace{{Name: "presence-v2", Projects: []model.Project{{Name: "baloo-service"}, {Name: "pistachio-api"}, {Name: "api"}}}}}}
	m.query = "api"
	m.filter()
	if len(m.rows) != 2 || m.rows[0].label != "api" || m.selected().label != "api" {
		t.Fatalf("rows: %+v", m.rows)
	}
	// "pre" only matches the workspace name, which a scoped picker ignores.
	m.query = "pre"
	m.filter()
	if len(m.rows) != 0 {
		t.Fatalf("workspace name matched: %+v", m.rows)
	}
}

func TestRefreshKeepsSelectionWhileFiltered(t *testing.T) {
	m := treeFixture()
	m.query = "w"
	m.filter()
	m.cursor = 1
	m.filter()
	if m.cursor != 1 {
		t.Fatal("refresh reset the selection")
	}
}
