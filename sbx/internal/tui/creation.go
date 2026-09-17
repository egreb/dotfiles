package tui

import (
	"fmt"
	"io"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type CreationStep struct {
	Name string
	Run  func() error
}

type creationDone struct{ err error }
type creationTick time.Time

type creation struct {
	workspace string
	steps     []CreationStep
	current   int
	frame     int
	err       error
}

func (m creation) next() tea.Cmd {
	return func() tea.Msg { return creationDone{m.steps[m.current].Run()} }
}
func creationPulse() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return creationTick(t) })
}
func (m creation) Init() tea.Cmd { return tea.Batch(m.next(), creationPulse()) }
func (m creation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case creationDone:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.current++
		if m.current == len(m.steps) {
			return m, tea.Quit
		}
		return m, m.next()
	case creationTick:
		m.frame++
		return m, creationPulse()
	}
	return m, nil
}
func (m creation) View() tea.View {
	var view strings.Builder
	fmt.Fprintf(&view, "Creating %s\n\nCloning projects (%d/%d complete)\n\n", m.workspace, m.current, len(m.steps))
	for i, step := range m.steps {
		status := "queued"
		if i < m.current {
			status = "complete"
		}
		if i == m.current {
			status = string("|/-\\"[m.frame%4]) + " cloning"
			if m.err != nil {
				status = "failed"
			}
		}
		fmt.Fprintf(&view, "  %s  %s\n", step.Name, status)
	}
	if m.err != nil {
		fmt.Fprintf(&view, "\n%s\n", m.err)
	}
	return tea.NewView(view.String())
}

// RunCreation leaves a completed project list in the terminal before setup continues.
// Input stays with subprocesses; cancellation comes from the caller's context.
func RunCreation(out io.Writer, workspace string, steps []CreationStep) error {
	if len(steps) == 0 {
		return nil
	}
	p := tea.NewProgram(creation{workspace: workspace, steps: steps}, tea.WithInput(nil), tea.WithOutput(out), tea.WithoutSignalHandler())
	result, err := p.Run()
	if err != nil {
		return err
	}
	return result.(creation).err
}
