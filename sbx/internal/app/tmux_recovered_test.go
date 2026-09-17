package app

import (
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
