package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sbx/internal/config"
	"sbx/internal/process"
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	for _, name := range []string{"alpha", "a.b", "a-b", "A2"} {
		if err := validateName(name); err != nil {
			t.Errorf("%q unexpectedly invalid: %v", name, err)
		}
	}
	for _, name := range []string{"", "a", "default", "_alpha", "alpha_beta", "../alpha"} {
		if err := validateName(name); err == nil {
			t.Errorf("%q unexpectedly valid", name)
		}
	}
}

func TestMarkerMatchesExactWorkspaceName(t *testing.T) {
	directory := t.TempDir()
	if err := writeMarker(directory, "alpha"); err != nil {
		t.Fatal(err)
	}
	if !markerMatches(directory, "alpha") {
		t.Fatal("marker did not match")
	}
	if markerMatches(directory, "beta") {
		t.Fatal("marker matched another workspace")
	}
	if err := os.WriteFile(filepath.Join(directory, ".sbx-managed"), []byte("format=1\nname=alpha-extra\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if markerMatches(directory, "alpha") {
		t.Fatal("marker accepted a prefix match")
	}
}

func TestSanitizePreview(t *testing.T) {
	if got := sanitizePreview("ok\x1b[31m\nnext\x00"); got != "ok\nnext" {
		t.Fatalf("got %q", got)
	}
}

func TestProjectNodeVersionSources(t *testing.T) {
	for _, tc := range []struct {
		name, file, contents string
		want, bad            bool
	}{
		{"node-version", ".node-version", "22.14.0", true, false},
		{"nvmrc", ".nvmrc", "lts/jod", true, false},
		{"engines", "package.json", `{"engines":{"node":">=22.13"}}`, true, false},
		{"unversioned", "package.json", `{}`, false, false},
		{"empty pin", ".nvmrc", "", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, tc.file), []byte(tc.contents), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := projectDeclaresNodeVersion(dir)
			if got != tc.want || (err != nil) != tc.bad {
				t.Fatalf("got %v, %v", got, err)
			}
		})
	}
}

func TestBootstrapUsesFNMBeforeHookAndPnpm(t *testing.T) {
	for _, failInstall := range []bool{false, true} {
		t.Run(fmt.Sprint(failInstall), func(t *testing.T) {
			dir := t.TempDir()
			bin := filepath.Join(dir, "bin")
			os.Mkdir(bin, 0700)
			log := filepath.Join(dir, "calls")
			script := `#!/bin/sh
printf '%s\n' "$*" >> "$TEST_FNM_LOG"
[ "$FNM_VERSION_FILE_STRATEGY" = local ] || exit 8
[ "$FNM_RESOLVE_ENGINES" = true ] || exit 8
if [ "$1" = install ]; then
 [ "$TEST_FNM_FAIL" != 1 ] || exit 9
 exit 0
fi
[ "$1" = exec ] && [ "$2" = -- ] || exit 10
shift 2
export TEST_NODE_SELECTED=1
exec "$@"
`
			os.WriteFile(filepath.Join(bin, "fnm"), []byte(script), 0700)
			os.WriteFile(filepath.Join(bin, "pnpm"), []byte("#!/bin/sh\n[ \"$TEST_NODE_SELECTED\" = 1 ] || exit 11\nprintf 'pnpm ran\\n' >> \"$TEST_FNM_LOG\"\n"), 0700)
			os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("22.14.0"), 0600)
			os.WriteFile(filepath.Join(dir, "pnpm-lock.yaml"), nil, 0600)
			t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("TEST_FNM_LOG", log)
			if failInstall {
				t.Setenv("TEST_FNM_FAIL", "1")
			} else {
				t.Setenv("TEST_FNM_FAIL", "0")
			}
			app := &App{Err: io.Discard, config: config.Config{HostShell: "/bin/sh"}, runner: process.Runner{Out: io.Discard, Err: io.Discard}}
			err := app.bootstrapProject(context.Background(), config.Project{Name: "frontend", BeforeBootstrap: `[ "$TEST_NODE_SELECTED" = 1 ]`}, dir)
			data, _ := os.ReadFile(log)
			calls := string(data)
			if failInstall {
				if err == nil || calls != "install\n" {
					t.Fatalf("setup continued after failure: %v %q", err, calls)
				}
			} else if err != nil || !strings.HasPrefix(calls, "install\nexec -- /bin/sh") || !strings.Contains(calls, "exec -- pnpm install\npnpm ran") {
				t.Fatalf("wrong Node context/order: %v %q", err, calls)
			}
		})
	}
}
