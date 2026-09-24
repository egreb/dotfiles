package app

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"sbx/internal/config"
	"sbx/internal/process"
)

func TestGitCacheRepairsMissingTreeOnOtherBranch(t *testing.T) {
	for _, brokenUpstream := range []bool{false, true} {
		name := "repair"
		if brokenUpstream {
			name = "failed-rebuild-preserves-cache"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			git := func(dir string, args ...string) string {
				t.Helper()
				cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v\n%s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			upstream := filepath.Join(root, "upstream")
			if err := os.Mkdir(upstream, 0755); err != nil {
				t.Fatal(err)
			}
			git(upstream, "init", "-b", "main")
			git(upstream, "config", "user.email", "test@example.com")
			git(upstream, "config", "user.name", "Test")
			git(upstream, "commit", "--allow-empty", "-m", "main")
			git(upstream, "checkout", "-b", "other")
			if err := os.MkdirAll(filepath.Join(upstream, "client", "database"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(upstream, "client", "database", "data"), []byte("branch-only data"), 0600); err != nil {
				t.Fatal(err)
			}
			git(upstream, "add", ".")
			git(upstream, "commit", "-m", "other branch tree")
			tree := git(upstream, "rev-parse", "HEAD:client/database")
			git(upstream, "checkout", "main")
			cacheRoot := filepath.Join(root, "cache")
			if err := os.Mkdir(cacheRoot, 0755); err != nil {
				t.Fatal(err)
			}
			cache := filepath.Join(cacheRoot, "project.git")
			// A local clone gives us loose objects we can deliberately remove.
			git(root, "clone", "--mirror", "--no-hardlinks", upstream, cache)
			object := filepath.Join("objects", tree[:2], tree[2:])
			if err := os.Remove(filepath.Join(cache, object)); err != nil {
				t.Fatal(err)
			}
			if brokenUpstream {
				if err := os.Remove(filepath.Join(upstream, ".git", object)); err != nil {
					t.Fatal(err)
				}
			}
			app := &App{Out: io.Discard, Err: io.Discard, runner: process.Runner{Out: io.Discard, Err: io.Discard}, config: config.Config{GitCacheRoot: cacheRoot}}
			p := config.Project{Name: "project", Repo: upstream, Branch: "main"}
			dest := filepath.Join(root, "workspace")
			err := app.cloneProjectFromCache(context.Background(), p, dest)
			if brokenUpstream {
				if err == nil {
					t.Fatal("accepted corrupt upstream")
				}
				if _, err := os.Stat(dest); !os.IsNotExist(err) {
					t.Fatalf("workspace created: %v", err)
				}
				git(cache, "rev-parse", "refs/heads/other")
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			git(cache, "fsck", "--full")
			git(dest, "fsck", "--full")
			git(dest, "checkout", "other")
			if got := git(dest, "remote", "get-url", "origin"); got != upstream {
				t.Fatalf("origin = %s", got)
			}
			backups, err := filepath.Glob(filepath.Join(cacheRoot, ".project-damaged-*", "mirror.git"))
			if err != nil || len(backups) != 1 {
				t.Fatalf("backups = %v, %v", backups, err)
			}
			// A healthy cache can be reused, with no shared pack inodes.
			second := filepath.Join(root, "second")
			if err := app.cloneProjectFromCache(context.Background(), p, second); err != nil {
				t.Fatal(err)
			}
			packs, _ := filepath.Glob(filepath.Join(second, ".git", "objects", "pack", "*.pack"))
			cachePacks, _ := filepath.Glob(filepath.Join(cache, "objects", "pack", "*.pack"))
			if len(packs) == 0 || len(cachePacks) == 0 {
				t.Fatal("expected packed transport clones")
			}
			for _, pack := range packs {
				info, err := os.Stat(pack)
				if err != nil {
					t.Fatal(err)
				}
				for _, cached := range cachePacks {
					cacheInfo, err := os.Stat(cached)
					if err != nil {
						t.Fatal(err)
					}
					if os.SameFile(info, cacheInfo) {
						t.Fatal("workspace pack hardlinked to cache")
					}
				}
			}
		})
	}
}
