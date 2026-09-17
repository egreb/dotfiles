package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"sbx/internal/config"
)

func TestListNamesWithoutExternalCommands(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"zulu", "alpha", "unmanaged", "mismatch"} {
		path := filepath.Join(root, name)
		if err := os.Mkdir(path, 0700); err != nil {
			t.Fatal(err)
		}
		if name == "unmanaged" {
			continue
		}
		markerName := name
		if name == "mismatch" {
			markerName = "different"
		}
		if err := writeMarker(path, markerName); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", t.TempDir())
	var out, errors bytes.Buffer
	app := &App{Out: &out, Err: &errors, configLoaded: true, config: config.Config{WorkspaceRoot: root}}
	if code := app.Run([]string{"list", "--names"}); code != 0 {
		t.Fatalf("exit %d: %s", code, &errors)
	}
	if got := out.String(); got != "alpha\nzulu\n" {
		t.Fatalf("unexpected names: %q", got)
	}
}
