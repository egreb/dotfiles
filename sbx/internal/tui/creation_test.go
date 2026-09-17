package tui

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCreationProgressAndFailure(t *testing.T) {
	ran := 0
	m := creation{workspace: "test", steps: []CreationStep{
		{Name: "api", Run: func() error { ran++; return nil }},
		{Name: "web", Run: func() error { ran++; return errors.New("clone failed") }},
	}}
	if view := m.View().Content; !strings.Contains(view, "api  | cloning") || !strings.Contains(view, "web  queued") {
		t.Fatal(view)
	}
	next, cmd := m.Update(m.next()())
	m = next.(creation)
	if view := m.View().Content; !strings.Contains(view, "api  complete") || !strings.Contains(view, "1/2 complete") {
		t.Fatal(view)
	}
	next, _ = m.Update(cmd())
	m = next.(creation)
	if ran != 2 || m.err == nil || !strings.Contains(m.View().Content, "web  failed") {
		t.Fatalf("unexpected state: %#v", m)
	}
}

func TestRunCreationCompletesSteps(t *testing.T) {
	var output bytes.Buffer
	count := 0
	err := RunCreation(&output, "example", []CreationStep{
		{Name: "api", Run: func() error { count++; return nil }},
		{Name: "web", Run: func() error { count++; return nil }},
	})
	if err != nil || count != 2 {
		t.Fatalf("count=%d err=%v", count, err)
	}
}
