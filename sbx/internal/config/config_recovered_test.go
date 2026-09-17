package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadProjects(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "projects.yaml")
	contents := `projects:
  - name: backend
    repo: git@example.test:backend.git
    before_bootstrap: make generate
  - name: frontend
    repo: git@example.test:frontend.git
    branch: develop
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	projects, err := LoadProjects(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 2 {
		t.Fatalf("got %d projects", len(projects))
	}
	if projects[0].Branch != "main" || projects[1].Branch != "develop" {
		t.Fatalf("unexpected branches: %#v", projects)
	}
}

func TestLoadProjectsRejectsAliases(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "projects.yaml")
	contents := "projects:\n  - &project\n    name: backend\n    repo: ok://backend\n  - *project\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadProjects(path)
	if err == nil || !strings.Contains(err.Error(), "aliases are not allowed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExpandHomeDoesNotEvaluateShellSyntax(t *testing.T) {
	home := "/home/tester"
	if got := expandHome("~/sandboxes", home); got != "/home/tester/sandboxes" {
		t.Fatalf("got %q", got)
	}
	if got := expandHome("$HOME/sandboxes", home); got != "$HOME/sandboxes" {
		t.Fatalf("got %q", got)
	}
}

func TestLoadUsesTmuxDefaultConfigWhenUnset(t *testing.T) {
	directory := t.TempDir()
	configDirectory := filepath.Join(directory, "config")
	if err := os.MkdirAll(configDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDirectory, "config.yaml"), []byte("host_shell: /bin/sh\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", directory)
	t.Setenv("SBX_CONFIG_DIR", configDirectory)
	t.Setenv("SBX_CONFIG_FILE", "")

	loaded, err := Load("/unused")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.TmuxConfig != "" {
		t.Fatalf("expected tmux's default config lookup, got %q", loaded.TmuxConfig)
	}
}
