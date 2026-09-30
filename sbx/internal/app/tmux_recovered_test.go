package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sbx/internal/config"
	"strings"
	"testing"
)

func TestTmuxArgsUseDefaultUserConfig(t *testing.T) {
	application := &App{config: config.Config{TmuxSocket: "sbx-test"}}
	want := []string{"-L", "sbx-test", "list-sessions"}
	if got := application.tmuxArgs("list-sessions"); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestNestedTmuxCommandUnsetsOuterContextAndQuotesTarget(t *testing.T) {
	command := nestedTmuxCommand([]string{"-L", "socket name", "attach-session", "-t", "alpha--front end"})
	if !strings.HasPrefix(command, "exec env -u TMUX -u TMUX_PANE tmux ") {
		t.Fatalf("outer tmux context was not removed: %q", command)
	}
	if !strings.Contains(command, "'socket name'") || !strings.Contains(command, "'alpha--front end'") {
		t.Fatalf("tmux arguments were not safely quoted: %q", command)
	}
}

func TestPrefixBindingCommandPreservesTmuxCommand(t *testing.T) {
	output := `bind-key -r -T prefix 1 select-window -t :=1`
	command, ok := prefixBindingCommand(output, "1")
	if !ok || command != "select-window -t :=1" {
		t.Fatalf("got %q, %v", command, ok)
	}
}

func TestTmuxArgsAllowExplicitUserConfig(t *testing.T) {
	application := &App{config: config.Config{TmuxSocket: "sbx-test", TmuxConfig: "/tmp/custom.conf"}}
	want := []string{"-L", "sbx-test", "-f", "/tmp/custom.conf", "list-sessions"}
	if got := application.tmuxArgs("list-sessions"); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAgentSessionReturnsToShellAfterExit(t *testing.T) {
	for _, status := range []int{0, 1, 130} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			directory := t.TempDir()
			native := filepath.Join(directory, "native cli")
			shell := filepath.Join(directory, "host shell")
			if err := os.WriteFile(native, []byte(fmt.Sprintf("#!/bin/sh\nprintf 'agent: <%%s>\\n' \"$@\"\nexit %d\n", status)), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(shell, []byte("#!/bin/sh\nprintf 'shell: <%s>\\n' \"$@\"\n"), 0755); err != nil {
				t.Fatal(err)
			}
			prompt := "fix 'quotes'; $(exit 99)"
			command := agentSessionCommand(native, shell, "example", prompt)
			output, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()
			if err != nil {
				t.Fatalf("shell did not survive agent exit: %v: %s", err, output)
			}
			want := "agent: <run>\nagent: <--name>\nagent: <example>\nagent: <-->\nagent: <" + prompt + ">\nshell: <-l>\n"
			if string(output) != want {
				t.Fatalf("got %q, want %q", output, want)
			}
		})
	}
}
