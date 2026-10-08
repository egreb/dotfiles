package app

import (
	"io"
	"os"
	"testing"

	"sbx/internal/config"
)

func TestNewDraftRoundTrip(t *testing.T) {
	app := &App{Out: io.Discard, Err: io.Discard, config: config.Config{ConfigDir: t.TempDir()}}
	if _, ok := app.loadNewDraft(); ok {
		t.Fatal("loaded a draft that was never saved")
	}
	if err := app.saveNewDraft("fix", "codex", "Fix the\nfailing tests", map[string]bool{"web": true, "api": true, "docs": false}); err != nil {
		t.Fatal(err)
	}
	draft, ok := app.loadNewDraft()
	if !ok || draft.Name != "fix" || draft.Agent != "codex" || draft.Prompt != "Fix the\nfailing tests" {
		t.Fatalf("unexpected draft: %#v", draft)
	}
	if len(draft.Selected) != 2 || draft.Selected[0] != "api" || draft.Selected[1] != "web" {
		t.Fatalf("unexpected selection: %#v", draft.Selected)
	}
	app.clearNewDraft()
	if _, err := os.Stat(app.newDraftPath()); !os.IsNotExist(err) {
		t.Fatal("draft was not cleared")
	}
}
