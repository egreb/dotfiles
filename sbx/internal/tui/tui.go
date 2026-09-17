package tui

import (
	"context"
	"io"
	"sbx/internal/model"
	"strings"
	"time"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

type Backend interface {
	Inventory(context.Context) (model.Inventory, error)
	Preview(context.Context, string, int) (string, error)
}
type ActionKind int

const (
	ActionNone ActionKind = iota
	ActionActivate
	ActionNew
	ActionDelete
)

type Action struct {
	Kind      ActionKind
	Selection model.Selection
}
type Options struct{ Workspace, SelectedSession, ExcludedSession, Title string }
type row struct {
	label, status, session string
	selection              model.Selection
}
type dashboard struct {
	backend               Backend
	options               Options
	inventory             model.Inventory
	rows                  []row
	cursor, width, height int
	query, preview, err   string
	showPane, showHelp    bool
	action                Action
}
type inventoryMsg struct {
	inventory model.Inventory
	err       error
}
type previewMsg struct {
	session, content string
	err              error
}
type tickMsg time.Time

func Run(b Backend, in io.Reader, out io.Writer, opts Options) (Action, error) {
	m := dashboard{backend: b, options: opts, width: 100, height: 30, showPane: true}
	v, e := tea.NewProgram(m, tea.WithInput(in), tea.WithOutput(out)).Run()
	if e != nil {
		return Action{}, e
	}
	return v.(dashboard).action, nil
}
func (m dashboard) load() tea.Cmd {
	return func() tea.Msg { v, e := m.backend.Inventory(context.Background()); return inventoryMsg{v, e} }
}
func pulse() tea.Cmd              { return tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }) }
func (m dashboard) Init() tea.Cmd { return tea.Batch(m.load(), pulse()) }
func (m dashboard) selected() row {
	if len(m.rows) == 0 {
		return row{}
	}
	return m.rows[min(m.cursor, len(m.rows)-1)]
}
func (m dashboard) loadPreview() tea.Cmd {
	s := m.selected().session
	h := m.height
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		v, e := m.backend.Preview(ctx, s, h)
		return previewMsg{s, v, e}
	}
}
func fuzzyMatch(s, q string) bool {
	r := []rune(strings.ToLower(q))
	i := 0
	for _, c := range strings.ToLower(s) {
		if i < len(r) && c == r[i] {
			i++
		}
	}
	return i == len(r)
}
func status(s *model.Session) string {
	if s == nil {
		return "unopened"
	}
	if s.PaneDead {
		return "exited"
	}
	if s.Attached {
		return "attached"
	}
	return s.CurrentCommand
}
func (m *dashboard) filter() {
	previous := m.selected()
	var rows []row
	for _, w := range m.inventory.Workspaces {
		if m.options.Workspace != "" && w.Name != m.options.Workspace {
			continue
		}
		add := func(name string, kind model.SelectionKind, s *model.Session) {
			label := name
			if m.options.Workspace == "" {
				label = w.Name + " :: " + name
			}
			session := ""
			if s != nil {
				session = s.Name
			}
			if session != "" && session == m.options.ExcludedSession {
				return
			}
			if !fuzzyMatch(label, m.query) {
				return
			}
			project := name
			if kind == model.SelectAgent {
				project = ""
			}
			rows = append(rows, row{label, status(s), session, model.Selection{Kind: kind, Workspace: w.Name, Project: project}})
		}
		add("agent", model.SelectAgent, w.AgentSession)
		for _, p := range w.Projects {
			add(p.Name, model.SelectProject, p.Session)
		}
	}
	m.rows = rows
	m.cursor = 0
	for i, r := range rows {
		if r.selection == previous.selection || m.options.SelectedSession != "" && r.session == m.options.SelectedSession {
			m.cursor = i
		}
	}
	m.options.SelectedSession = ""
}
func (m dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = v.Width
		m.height = v.Height
	case tickMsg:
		return m, tea.Batch(m.load(), pulse())
	case inventoryMsg:
		if v.err != nil {
			m.err = v.err.Error()
		} else {
			m.err = ""
			m.inventory = v.inventory
			m.filter()
		}
		return m, m.loadPreview()
	case previewMsg:
		if v.session == m.selected().session {
			m.preview = v.content
			if v.err != nil {
				m.preview = v.err.Error()
			}
		}
	case tea.PasteMsg:
		m.query += cleanFormText(v.Content, false)
		m.filter()
		return m, m.loadPreview()
	case tea.KeyPressMsg:
		switch v.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.showHelp {
				m.showHelp = false
			} else if m.query != "" {
				m.query = ""
				m.filter()
				return m, m.loadPreview()
			} else {
				return m, tea.Quit
			}
		case "enter":
			if len(m.rows) > 0 {
				m.action = Action{ActionActivate, m.selected().selection}
				return m, tea.Quit
			}
		case "ctrl+d":
			if len(m.rows) > 0 {
				m.action = Action{ActionDelete, m.selected().selection}
				return m, tea.Quit
			}
		case "ctrl+j", "down":
			m.cursor = min(max(0, len(m.rows)-1), m.cursor+1)
			return m, m.loadPreview()
		case "ctrl+k", "up":
			m.cursor = max(0, m.cursor-1)
			return m, m.loadPreview()
		case "ctrl+h":
			m.showPane = false
		case "ctrl+l":
			m.showPane = true
		case "tab":
			m.showPane = !m.showPane
		case "f1", "ctrl+g":
			m.showHelp = !m.showHelp
		case "ctrl+r":
			return m, m.load()
		case "backspace":
			r := []rune(m.query)
			if len(r) > 0 {
				m.query = string(r[:len(r)-1])
				m.filter()
				return m, m.loadPreview()
			}
		default:
			if v.Text != "" {
				for _, r := range v.Text {
					if unicode.IsPrint(r) {
						m.query += string(r)
					}
				}
				m.filter()
				return m, m.loadPreview()
			}
		}
	}
	return m, nil
}
func fit(s string, w int) string {
	w = max(0, w)
	s = ansi.Truncate(s, w, "…")
	return s + strings.Repeat(" ", max(0, w-ansi.StringWidth(s)))
}
func (m dashboard) View() tea.View {
	w := max(1, m.width-2)
	var lines []string
	title := m.options.Title
	if title == "" {
		title = "sbx sessions"
	}
	lines = append(lines, title, "Filter: "+m.query+"█")
	// Full selected identity occupies its own wrapped lines, independent of preview width.
	selected := m.selected()
	identity := "Select a workspace or project"
	if len(m.rows) > 0 {
		identity = "Open: " + selected.selection.Workspace + " / "
		if selected.selection.Kind == model.SelectAgent {
			identity += "agent"
		} else {
			identity += selected.selection.Project
		}
	}
	lines = append(lines, strings.Split(ansi.Wrap(identity, w, ""), "\n")...)
	if m.err != "" {
		lines = append(lines, "Error: "+m.err)
	}
	if m.showHelp {
		lines = append(lines, "Type: filter   Ctrl-j/k: select   Enter: open", "Tab / Ctrl-h/l: preview   Ctrl-r: refresh", "Ctrl-d: delete (confirmation follows)   Esc: clear/quit")
	}
	h := max(1, m.height-len(lines)-3)
	left := w
	show := m.showPane && w >= 90
	if show {
		left = w * 3 / 5
	}
	preview := strings.Split(m.preview, "\n")
	offset := max(0, m.cursor-h+1)
	for i := 0; i < h; i++ {
		a := ""
		if offset+i < len(m.rows) {
			r := m.rows[offset+i]
			a = fit(r.label, max(0, left-13)) + " " + fit(r.status, 11)
			a = fit(a, left)
			if offset+i == m.cursor {
				a = "\x1b[7m" + a + "\x1b[0m"
			}
		} else {
			a = fit("", left)
		}
		if show {
			b := ""
			if i < len(preview) {
				b = preview[i]
			}
			a += " │ " + fit(b, w-left-3)
		}
		lines = append(lines, a)
	}
	lines = append(lines, "Ctrl-j/k select · Enter open · Tab preview · Esc clear/quit · F1 help")
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], w, "…")
	}
	if len(lines) > m.height {
		lines = lines[:m.height]
	}
	v := tea.NewView(strings.Join(lines, "\n"))
	v.AltScreen = true
	return v
}
