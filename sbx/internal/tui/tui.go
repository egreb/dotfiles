package tui

import (
	"context"
	"fmt"
	"io"
	"sbx/internal/model"
	"sort"
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
	ActionDeleteConfirmed
)

type Action struct {
	Kind      ActionKind
	Selection model.Selection
}
type Options struct {
	Workspace, SelectedSession, ExcludedSession, Title string
	DeleteMode                                         bool
}
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
	searching, sortRecent bool
	collapsed             map[string]bool
	remembered            map[string]model.Selection
	pendingDelete         string
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
	m := dashboard{backend: b, options: opts, width: 100, height: 30, sortRecent: true}
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
	if !m.showPane {
		return nil
	}
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
	workspaces := append([]model.Workspace(nil), m.inventory.Workspaces...)
	sort.SliceStable(workspaces, func(i, j int) bool {
		a, b := workspaces[i], workspaces[j]
		if m.sortRecent && a.LastUsed != b.LastUsed {
			return a.LastUsed > b.LastUsed
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	for _, w := range workspaces {
		if m.options.Workspace != "" && w.Name != m.options.Workspace {
			continue
		}
		if m.options.DeleteMode {
			if fuzzyMatch(w.Name, m.query) {
				session := ""
				if w.AgentSession != nil {
					session = w.AgentSession.Name
				}
				rows = append(rows, row{w.Name, status(w.AgentSession), session, model.Selection{Kind: model.SelectWorkspace, Workspace: w.Name}})
			}
			continue
		}
		add := func(name string, kind model.SelectionKind, s *model.Session) {
			label := name
			session := ""
			if s != nil {
				session = s.Name
			}
			if session != "" && session == m.options.ExcludedSession {
				return
			}
			if !fuzzyMatch(w.Name+" :: "+name, m.query) {
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
		if !m.searching {
			return m, nil
		}
		m.query += cleanFormText(v.Content, false)
		m.filter()
		return m, m.loadPreview()
	case tea.KeyPressMsg:
		if m.pendingDelete != "" {
			switch v.String() {
			case "y", "Y":
				m.action = Action{ActionDeleteConfirmed, model.Selection{Kind: model.SelectWorkspace, Workspace: m.pendingDelete}}
				return m, tea.Quit
			case "n", "N", "esc", "enter", "ctrl+c":
				m.pendingDelete = ""
			}
			return m, nil
		}
		if m.searching {
			switch v.String() {
			case "esc", "enter":
				m.searching = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			case "backspace":
				r := []rune(m.query)
				if len(r) > 0 {
					m.query = string(r[:len(r)-1])
				}
			default:
				for _, r := range v.Text {
					if unicode.IsPrint(r) {
						m.query += string(r)
					}
				}
			}
			m.filter()
			m.preview = ""
			return m, m.loadPreview()
		}
		switch v.String() {
		case "/":
			m.searching = true
			return m, nil
		case "s":
			m.sortRecent = !m.sortRecent
			m.filter()
			return m, nil
		case "h", "left":
			m.moveWorkspace(-1)
			m.preview = ""
			return m, m.loadPreview()
		case "l", "right":
			m.moveWorkspace(1)
			m.preview = ""
			return m, m.loadPreview()
		case "space":
			if m.collapsed == nil {
				m.collapsed = map[string]bool{}
			}
			name := m.selected().selection.Workspace
			m.collapsed[name] = !m.collapsed[name]
			return m, nil
		}
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
				kind := ActionActivate
				if m.options.DeleteMode {
					kind = ActionDelete
				}
				m.action = Action{kind, m.selected().selection}
				return m, tea.Quit
			}
		case "D", "ctrl+d":
			if len(m.rows) > 0 {
				if m.options.DeleteMode {
					m.action = Action{ActionDelete, m.selected().selection}
					return m, tea.Quit
				}
				m.pendingDelete = m.selected().selection.Workspace
			}
		case "j", "ctrl+j", "down":
			m.moveSession(1)
			m.preview = ""
			return m, m.loadPreview()
		case "k", "ctrl+k", "up":
			m.moveSession(-1)
			m.preview = ""
			return m, m.loadPreview()
		case "ctrl+h":
			m.showPane = false
		case "ctrl+l":
			m.showPane = true
			return m, m.loadPreview()
		case "tab":
			m.showPane = !m.showPane
			return m, m.loadPreview()
		case "f1", "ctrl+g":
			m.showHelp = !m.showHelp
		case "ctrl+r":
			return m, m.load()
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
	if m.pendingDelete != "" {
		return (confirmation{question: "Delete workspace " + m.pendingDelete + "?\n\nThis permanently removes its files and sandbox and stops all its sessions.\nUncommitted changes will be lost.", width: m.width}).View()
	}
	w := max(1, m.width-2)
	var lines []string
	title := m.options.Title
	if title == "" {
		title = "sbx sessions"
	}
	sortLabel := "name A–Z"
	if m.sortRecent {
		sortLabel = "last used"
	}
	filterLabel := "/ filter"
	if m.query != "" {
		filterLabel = "Filter: " + m.query
	}
	if m.searching {
		filterLabel = "Filter: " + m.query + "█  (Enter: done)"
	}
	lines = append(lines, title+" · sort: "+sortLabel+" (s)", filterLabel)
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
	if m.options.DeleteMode && len(m.rows) > 0 {
		identity = "Delete workspace: " + selected.selection.Workspace + " (confirmation follows)"
	}
	lines = append(lines, strings.Split(ansi.Wrap(identity, w, ""), "\n")...)
	if m.err != "" {
		lines = append(lines, "Error: "+m.err)
	}
	verb := "open"
	if m.options.DeleteMode {
		verb = "select for deletion"
	}
	if m.showHelp {
		lines = append(lines, "h/l: workspace   j/k: session   Enter: "+verb, "Tab / Ctrl-h/l: preview   Ctrl-r: refresh", "/: filter (Enter ends editing)   s: sort   Space: fold", "D / Ctrl-d: delete workspace (confirmation follows)   Esc: clear/quit")
	}
	h := max(1, m.height-len(lines)-3)
	left := w
	show := m.showPane && w >= 90
	if show {
		left = w * 3 / 5
	}
	preview := strings.Split(m.preview, "\n")
	display, active := m.treeLines(left)
	offset := max(0, active-h+1)
	for i := 0; i < h; i++ {
		a := ""
		if offset+i < len(display) {
			a = display[offset+i]
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
	lines = append(lines, "h/l workspace · j/k session · / filter · s sort · Space fold · Enter "+verb+" · D delete · Tab preview · F1 help")
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

// Workspace movement follows the visible sort order and restores its last selection.
func (m *dashboard) moveWorkspace(delta int) {
	if len(m.rows) == 0 {
		return
	}
	if m.remembered == nil {
		m.remembered = map[string]model.Selection{}
	}
	current := m.selected().selection
	m.remembered[current.Workspace] = current
	starts := []int{0}
	group := 0
	for i := 1; i < len(m.rows); i++ {
		if m.rows[i].selection.Workspace != m.rows[i-1].selection.Workspace {
			starts = append(starts, i)
		}
		if i == m.cursor {
			group = len(starts) - 1
		}
	}
	target := starts[max(0, min(len(starts)-1, group+delta))]
	name := m.rows[target].selection.Workspace
	m.cursor = target
	for i := target; i < len(m.rows) && m.rows[i].selection.Workspace == name; i++ {
		if m.rows[i].selection == m.remembered[name] {
			m.cursor = i
		}
	}
	if m.collapsed != nil {
		delete(m.collapsed, name)
	}
}
func (m *dashboard) moveSession(delta int) {
	if m.options.DeleteMode {
		m.moveWorkspace(delta)
		return
	}
	if m.collapsed != nil {
		delete(m.collapsed, m.selected().selection.Workspace)
	}
	next := m.cursor + delta
	if next >= 0 && next < len(m.rows) && m.rows[next].selection.Workspace == m.selected().selection.Workspace {
		m.cursor = next
	}
}
func (m dashboard) treeLines(width int) ([]string, int) {
	var lines []string
	active := 0
	selected := m.selected().selection.Workspace
	for i := 0; i < len(m.rows); {
		end := i + 1
		name := m.rows[i].selection.Workspace
		for end < len(m.rows) && m.rows[end].selection.Workspace == name {
			end++
		}
		open := name == selected && !m.collapsed[name] && !m.options.DeleteMode
		marker := "▸ "
		if open {
			marker = "▾ "
		}
		last := "last used unknown"
		for _, w := range m.inventory.Workspaces {
			if w.Name == name && w.LastUsed > 0 {
				last = time.Unix(w.LastUsed, 0).Format("Jan 02 15:04")
				break
			}
		}
		header := fit(fmt.Sprintf("%s%s  · %d sessions · %s", marker, name, end-i, last), width)
		if m.options.DeleteMode {
			header = fit(name+" · "+last, width)
		}
		if name == selected {
			active = len(lines)
			header = "\x1b[1m" + header + "\x1b[0m"
			if !open {
				header = "\x1b[7m" + header + "\x1b[0m"
			}
		}
		lines = append(lines, header)
		if open {
			for j := i; j < end; j++ {
				r := m.rows[j]
				line := fit("    "+r.label, max(0, width-13)) + " " + fit(r.status, 12)
				line = fit(line, width)
				if j == m.cursor {
					active = len(lines)
					line = "\x1b[7m" + line + "\x1b[0m"
				}
				lines = append(lines, line)
			}
		}
		i = end
	}
	if len(lines) == 0 {
		lines = append(lines, fit("No matching workspaces or sessions", width))
	}
	return lines, active
}
