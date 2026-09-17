package tui

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

type Repository struct{ Name, Branch string }
type WorkspaceForm struct {
	Name, Agent, Prompt string
	Repositories        []Repository
	Selected            map[string]bool
	ValidateName        func(string) error
}
type workspaceForm struct {
	data                             WorkspaceForm
	focus, repoCursor, width, height int
	nameCursor, promptCursor         int
	selectedOffset                   int
	query, err                       string
	accepted                         bool
}

func RunWorkspaceForm(in io.Reader, out io.Writer, data WorkspaceForm) (WorkspaceForm, bool, error) {
	m := newWorkspaceForm(data)
	result, err := tea.NewProgram(m, tea.WithInput(in), tea.WithOutput(out)).Run()
	if err != nil {
		return data, false, err
	}
	final := result.(workspaceForm)
	return final.data, final.accepted, nil
}
func newWorkspaceForm(data WorkspaceForm) workspaceForm {
	selected := make(map[string]bool)
	for name, value := range data.Selected {
		selected[name] = value
	}
	data.Selected = selected
	return workspaceForm{data: data, width: 90, height: 32, nameCursor: len([]rune(data.Name)), promptCursor: len([]rune(data.Prompt))}
}
func (m workspaceForm) Init() tea.Cmd { return nil }
func (m workspaceForm) matches() []Repository {
	var rows []Repository
	for _, repo := range m.data.Repositories {
		if strings.Contains(strings.ToLower(repo.Name), strings.ToLower(m.query)) {
			rows = append(rows, repo)
		}
	}
	return rows
}
func (m workspaceForm) selectedRepositories() []Repository {
	var selected []Repository
	for _, repo := range m.data.Repositories {
		if m.data.Selected[repo.Name] {
			selected = append(selected, repo)
		}
	}
	return selected
}
func (m workspaceForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.PasteMsg:
		if m.focus == 0 {
			insertText(&m.data.Name, &m.nameCursor, cleanFormText(msg.Content, false))
		}
		if m.focus == 2 {
			insertText(&m.data.Prompt, &m.promptCursor, cleanFormText(msg.Content, true))
		}
	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
		case "pgup":
			m.selectedOffset = max(0, m.selectedOffset-1)
			return m, nil
		case "pgdown":
			m.selectedOffset = min(max(0, len(m.selectedRepositories())-1), m.selectedOffset+1)
			return m, nil
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % 5
			return m, nil
		case "shift+tab":
			m.focus = (m.focus + 4) % 5
			return m, nil
		}
		m.err = ""
		switch m.focus {
		case 0:
			editFormText(&m.data.Name, &m.nameCursor, msg, false)
		case 1:
			if key == "left" || key == "right" || key == " " || key == "space" {
				if m.data.Agent == "claude" {
					m.data.Agent = "codex"
				} else {
					m.data.Agent = "claude"
				}
			}
		case 2:
			editFormText(&m.data.Prompt, &m.promptCursor, msg, true)
		case 3:
			rows := m.matches()
			switch key {
			case "up", "ctrl+k":
				m.repoCursor = max(0, m.repoCursor-1)
			case "down", "ctrl+j":
				m.repoCursor = min(max(0, len(rows)-1), m.repoCursor+1)
			case " ", "space":
				if len(rows) > 0 {
					name := rows[m.repoCursor].Name
					m.data.Selected[name] = !m.data.Selected[name]
				}
			case "backspace":
				r := []rune(m.query)
				if len(r) > 0 {
					m.query = string(r[:len(r)-1])
				}
				m.repoCursor = 0
			default:
				if msg.Text != "" {
					m.query += cleanFormText(msg.Text, false)
					m.repoCursor = 0
				}
			}
		case 4:
			if key == "enter" || key == " " || key == "space" {
				if m.data.ValidateName != nil {
					if err := m.data.ValidateName(m.data.Name); err != nil {
						m.err = err.Error()
						m.focus = 0
						return m, nil
					}
				}
				count := 0
				for _, repo := range m.data.Repositories {
					if m.data.Selected[repo.Name] {
						count++
					}
				}
				if count == 0 {
					m.err = "Select at least one repository."
					m.focus = 3
					return m, nil
				}
				m.accepted = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
func cleanFormText(s string, multiline bool) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' && multiline {
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
func insertText(value *string, cursor *int, text string) {
	r := []rune(*value)
	added := []rune(text)
	r = append(r[:*cursor], append(added, r[*cursor:]...)...)
	*cursor += len(added)
	*value = string(r)
}
func editFormText(value *string, cursor *int, key tea.KeyPressMsg, multiline bool) {
	r := []rune(*value)
	switch key.String() {
	case "left":
		*cursor = max(0, *cursor-1)
	case "right":
		*cursor = min(len(r), *cursor+1)
	case "home", "ctrl+a":
		*cursor = 0
	case "end", "ctrl+e":
		*cursor = len(r)
	case "backspace":
		if *cursor > 0 {
			*value = string(append(r[:*cursor-1], r[*cursor:]...))
			*cursor--
		}
	case "delete":
		if *cursor < len(r) {
			*value = string(append(r[:*cursor], r[*cursor+1:]...))
		}
	case "enter":
		if multiline {
			insertText(value, cursor, "\n")
		}
	default:
		if key.Text != "" {
			insertText(value, cursor, cleanFormText(key.Text, multiline))
		}
	}
}
func (m workspaceForm) View() tea.View {
	width := max(20, min(m.width-4, 100))
	label := func(index int, text string) string {
		if m.focus == index {
			return "\x1b[1;36m› " + text + "\x1b[0m"
		}
		return "  " + text
	}
	field := func(value string, cursor, index int, placeholder string) string {
		r := []rune(value)
		if m.focus == index {
			value = string(r[:cursor]) + "▏" + string(r[cursor:])
		}
		if value == "" {
			value = "\x1b[2m" + placeholder + "\x1b[0m"
		}
		return value
	}
	var b strings.Builder
	b.WriteString("\n  \x1b[1mNew workspace\x1b[0m\n  Choose a task and the repositories it needs.\n\n")
	fmt.Fprintf(&b, "%s\n    %s\n\n", label(0, "Name"), ansi.Truncate(field(m.data.Name, m.nameCursor, 0, "e.g. fix-checkout"), width-4, "…"))
	agents := "● Claude     ○ Codex"
	if m.data.Agent == "codex" {
		agents = "○ Claude     ● Codex"
	}
	fmt.Fprintf(&b, "%s  %s\n\n", label(1, "Agent"), agents)
	fmt.Fprintf(&b, "%s \x1b[2m(optional)\x1b[0m\n", label(2, "Task prompt"))
	prompt := field(m.data.Prompt, m.promptCursor, 2, "What should the agent work on?")
	lines := strings.Split(ansi.Hardwrap(prompt, max(10, width-4), true), "\n")
	// Keep the insertion point visible for long prompts.
	start := 0
	if m.focus == 2 {
		for i, line := range lines {
			if strings.Contains(line, "▏") {
				start = max(0, i-2)
				break
			}
		}
	}
	for i := start; i < min(len(lines), start+3); i++ {
		fmt.Fprintf(&b, "    %s\n", lines[i])
	}
	for i := min(len(lines)-start, 3); i < 3; i++ {
		b.WriteString("\n")
	}
	count := 0
	for _, repo := range m.data.Repositories {
		if m.data.Selected[repo.Name] {
			count++
		}
	}
	fmt.Fprintf(&b, "\n%s  %d/%d selected to clone\n", label(3, "Repositories"), count, len(m.data.Repositories))
	filter := m.query
	if filter == "" {
		filter = "type to filter"
	}
	fmt.Fprintf(&b, "  \x1b[2m%s · ↑/↓ move · Space select\x1b[0m\n", filter)
	rows := m.matches()
	visible := max(1, m.height-26)
	visible = min(visible, 12)
	offset := max(0, m.repoCursor-visible+1)
	selected := m.selectedRepositories()
	selectedStart := min(m.selectedOffset, max(0, len(selected)-visible))
	columnWidth := max(8, (width-6)/2)
	leftHeader := "Available repositories"
	rightHeader := fmt.Sprintf("Will clone (%d)", len(selected))
	fmt.Fprintf(&b, "   %-*s │ %s\n", columnWidth, leftHeader, rightHeader)
	for lineIndex := 0; lineIndex < min(visible, max(1, max(len(rows)-offset, len(selected)-selectedStart))); lineIndex++ {
		left := ""
		i := offset + lineIndex
		if i < len(rows) {
			repo := rows[i]
			check := "[ ]"
			if m.data.Selected[repo.Name] {
				check = "[✓]"
			}
			left = ansi.Truncate(check+" "+repo.Name, columnWidth, "…")
		} else if lineIndex == 0 {
			left = "No matches"
		}
		left += strings.Repeat(" ", max(0, columnWidth-ansi.StringWidth(left)))
		if m.focus == 3 && i == m.repoCursor && i < len(rows) {
			left = "\x1b[7m" + left + "\x1b[0m"
		}
		right := ""
		if selectedStart+lineIndex < len(selected) {
			right = selected[selectedStart+lineIndex].Name
		} else if len(selected) == 0 && lineIndex == 0 {
			right = "None selected"
		}
		fmt.Fprintf(&b, "   %s │ %s\n", left, ansi.Truncate(right, columnWidth, "…"))
	}
	fmt.Fprintf(&b, "  \x1b[2m%d matching · PgUp/PgDn scroll selected\x1b[0m\n", len(rows))
	fmt.Fprintf(&b, "\n%s\n", label(4, "[ Create workspace ]"))
	if m.err != "" {
		fmt.Fprintf(&b, "  \x1b[31m%s\x1b[0m\n", ansi.Truncate(m.err, width, "…"))
	}
	b.WriteString("\n  \x1b[2mTab / Shift-Tab move · ←/→ change agent · Esc cancel\x1b[0m")
	content := b.String()
	if m.height < 28 {
		content = strings.ReplaceAll(content, "\n\n", "\n")
	}
	contentLines := strings.Split(content, "\n")
	for i, line := range contentLines {
		contentLines[i] = ansi.Truncate(line, max(1, m.width), "…")
	}
	view := tea.NewView(strings.Join(contentLines, "\n"))
	view.AltScreen = true
	return view
}
