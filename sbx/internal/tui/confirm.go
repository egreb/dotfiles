package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"io"
)

type confirmation struct {
	question string
	width    int
	accepted bool
}

func (m confirmation) Init() tea.Cmd { return nil }
func (m confirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = v.Width
	case tea.KeyPressMsg:
		switch v.String() {
		case "y", "Y":
			m.accepted = true
			return m, tea.Quit
		case "n", "N", "esc", "enter", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m confirmation) View() tea.View {
	v := tea.NewView(ansi.Wrap(m.question+"\n\n[y] Yes   [n / Enter / Esc] No", max(1, m.width-2), ""))
	v.AltScreen = true
	return v
}
func Confirm(in io.Reader, out io.Writer, question string) (bool, error) {
	result, err := tea.NewProgram(confirmation{question: question, width: 80}, tea.WithInput(in), tea.WithOutput(out)).Run()
	if err != nil {
		return false, err
	}
	return result.(confirmation).accepted, nil
}
