package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/config"
	"sbx/internal/process"
	"sort"
	"strings"
)

func canonicalDirectory(path string) (string, error) {
	p, e := filepath.Abs(path)
	if e != nil {
		return "", e
	}
	p, e = filepath.EvalSymlinks(p)
	if e != nil {
		return "", e
	}
	i, e := os.Stat(p)
	if e != nil {
		return "", e
	}
	if !i.IsDir() {
		return "", fmt.Errorf("not a directory: %s", p)
	}
	return p, nil
}
func (app *App) canonicalWorkspaceRoot(create bool) (string, error) {
	if create {
		if e := os.MkdirAll(app.config.WorkspaceRoot, 0755); e != nil {
			return "", e
		}
	}
	return canonicalDirectory(app.config.WorkspaceRoot)
}
func writeMarker(path, name string) error {
	return os.WriteFile(filepath.Join(path, ".sbx-managed"), []byte("format=1\nname="+name+"\n"), 0600)
}
func markerMatches(path, name string) bool {
	f := filepath.Join(path, ".sbx-managed")
	i, e := os.Lstat(f)
	if e != nil || !i.Mode().IsRegular() {
		return false
	}
	b, e := os.ReadFile(f)
	if e != nil {
		return false
	}
	var format, stored string
	for _, l := range strings.Split(string(b), "\n") {
		k, v, _ := strings.Cut(l, "=")
		switch k {
		case "format":
			format = v
		case "name":
			stored = v
		}
	}
	return format == "1" && stored == name
}
func (app *App) workspaceCandidates() ([]string, error) {
	entries, e := os.ReadDir(app.config.WorkspaceRoot)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() || validateName(entry.Name()) != nil {
			continue
		}
		if _, e := app.assertSafeWorkspace(filepath.Join(app.config.WorkspaceRoot, entry.Name()), entry.Name()); e == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
func (app *App) updateProjectGitCache(ctx context.Context, p config.Project) error {
	if e := os.MkdirAll(app.config.GitCacheRoot, 0755); e != nil {
		return e
	}
	cache := filepath.Join(app.config.GitCacheRoot, p.Name+".git")
	if _, e := os.Stat(cache); e == nil {
		origin, e := app.runner.Output(ctx, "git", []string{"-C", cache, "remote", "get-url", "origin"}, process.Options{})
		if e != nil {
			return e
		}
		if strings.TrimSpace(origin) != p.Repo {
			return fmt.Errorf("cache origin mismatch for %s", p.Name)
		}
		return app.runner.Run(ctx, "git", []string{"-C", cache, "fetch", "--prune", "origin"}, process.Options{})
	} else if !os.IsNotExist(e) {
		return e
	}
	temp, e := os.MkdirTemp(app.config.GitCacheRoot, "."+p.Name+"-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(temp)
	if e = app.runner.Run(ctx, "git", []string{"clone", "--mirror", "--", p.Repo, temp}, process.Options{}); e != nil {
		return e
	}
	return os.Rename(temp, cache)
}
