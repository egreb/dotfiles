package app

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestRemoveTreeWithRetry(t *testing.T) {
	for _, tc := range []struct {
		name                string
		failure             error
		failures, wantCalls int
		wantErr             bool
	}{
		{"success", nil, 0, 1, false},
		{"transient nonempty", syscall.ENOTEMPTY, 2, 3, false},
		{"persistent nonempty", syscall.ENOTEMPTY, 10, 6, true},
		{"permission denied", syscall.EACCES, 1, 1, true},
		{"io failure", syscall.EIO, 1, 1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls, waits := 0, 0
			err := removeTreeWithRetry("/workspace", func(path string) error {
				if path != "/workspace" {
					t.Fatalf("unexpected path: %s", path)
				}
				calls++
				if calls <= tc.failures {
					return &os.PathError{Op: "unlinkat", Path: path + "/project", Err: tc.failure}
				}
				return nil
			}, func(delay time.Duration) {
				if delay <= 0 || delay > 2*time.Second {
					t.Fatalf("unexpected delay: %s", delay)
				}
				waits++
			})
			if calls != tc.wantCalls || waits != calls-1 {
				t.Fatalf("calls=%d waits=%d", calls, waits)
			}
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr && !errors.Is(err, tc.failure) {
				t.Fatalf("lost underlying error: %v", err)
			}
		})
	}
}

func TestRemoveWorkspaceTree(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "keep")
	if err := os.WriteFile(outside, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	workspace := filepath.Join(root, "workspace")
	project := filepath.Join(workspace, "project")
	if err := os.MkdirAll(project, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(project, "link")); err != nil {
		t.Fatal(err)
	}
	if err := removeWorkspaceTree(workspace); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(workspace); !os.IsNotExist(err) {
		t.Fatalf("workspace still exists: %v", err)
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatalf("symlink target removed: %v", err)
	}
	if err := removeWorkspaceTree(workspace); err != nil {
		t.Fatalf("already absent: %v", err)
	}
}
