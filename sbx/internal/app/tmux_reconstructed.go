package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/model"
	"sbx/internal/process"
	"strconv"
	"strings"
)

func (app *App) tmuxOption(ctx context.Context, session, key string) string {
	v, _ := app.tmuxOutput(ctx, true, "show-options", "-qv", "-t", session, key)
	return trimCommandOutput(v)
}
func (app *App) tmuxSessionManaged(ctx context.Context, name string) bool {
	return app.tmuxSessionExists(ctx, name) && app.tmuxOption(ctx, name, "@sbx-managed") == "1" && app.tmuxOption(ctx, name, "@sbx-name") == name && app.tmuxOption(ctx, name, "@sbx-role") == "agent"
}
func (app *App) listTmuxSessions(ctx context.Context) []tmuxSession {
	format := strings.Join([]string{"#{session_name}", "#{@sbx-managed}", "#{@sbx-name}", "#{@sbx-workspace}", "#{@sbx-docker-name}", "#{@sbx-role}", "#{@sbx-project}", "#{@sbx-last-session}", "#{session_attached}", "#{pane_id}", "#{pane_current_command}", "#{pane_dead}"}, "\t")
	v, e := app.tmuxOutput(ctx, true, "list-sessions", "-F", format)
	if e != nil {
		return nil
	}
	var sessions []tmuxSession
	for _, l := range outputLines(v) {
		f := strings.Split(l, "\t")
		if len(f) != 12 {
			continue
		}
		attached, _ := strconv.Atoi(f[8])
		sessions = append(sessions, tmuxSession{Name: f[0], Managed: f[1] == "1", WorkspaceName: f[2], Workspace: f[3], DockerName: f[4], Role: f[5], Project: f[6], LastSession: f[7], Attached: attached > 0, ActivePaneID: f[9], CurrentCommand: f[10], PaneDead: f[11] == "1"})
	}
	return sessions
}
func (app *App) currentManagedContext(ctx context.Context) (managedContext, error) {
	// Popups are not panes; capture the invoking pane when the binding runs.
	pane := os.Getenv("SBX_TMUX_PANE")
	if pane == "" {
		pane = os.Getenv("TMUX_PANE")
	}
	if pane == "" {
		return managedContext{}, fmt.Errorf("this command must run inside a managed tmux pane")
	}
	out, e := app.tmuxOutput(ctx, true, "display-message", "-p", "-t", pane, "#{session_name}")
	if e != nil {
		return managedContext{}, e
	}
	session := trimCommandOutput(out)
	if app.tmuxOption(ctx, session, "@sbx-managed") != "1" {
		return managedContext{}, fmt.Errorf("not a managed session: %s", session)
	}
	c := managedContext{Session: session, Pane: pane, WorkspaceName: app.tmuxOption(ctx, session, "@sbx-name"), Workspace: app.tmuxOption(ctx, session, "@sbx-workspace"), DockerName: app.tmuxOption(ctx, session, "@sbx-docker-name")}
	if c.WorkspaceName == "" || c.DockerName != c.WorkspaceName {
		return c, fmt.Errorf("invalid workspace metadata")
	}
	_, e = app.assertSafeWorkspace(c.Workspace, c.WorkspaceName)
	return c, e
}
func (app *App) setWorkspaceLastSession(ctx context.Context, name, target string) error {
	return app.tmuxRun(ctx, "set-option", "-t", name, "@sbx-last-session", target)
}
func (app *App) switchCurrentClient(ctx context.Context, target string) error {
	args := []string{"switch-client", "-t", "=" + target}
	if client := os.Getenv("SBX_TMUX_CLIENT"); client != "" {
		args = append(args, "-c", client)
	}
	return app.tmuxRun(ctx, args...)
}
func (app *App) createAgentTmuxSession(ctx context.Context, name, workspace, last, prompt string) error {
	if e := validateName(name); e != nil {
		return e
	}
	if e := app.needNative(); e != nil {
		return e
	}
	if last == "" {
		last = name
	}
	command := "exec " + process.ShellQuote(app.native) + " run --name " + process.ShellQuote(name)
	if prompt != "" {
		command += " -- " + process.ShellQuote(prompt)
	}
	if e := app.tmuxRun(ctx, "new-session", "-d", "-s", name, "-n", "agent", "-c", workspace, command); e != nil {
		return e
	}
	configured := false
	defer func() {
		if !configured {
			_ = app.tmuxRunQuiet(context.Background(), "kill-session", "-t", "="+name)
		}
	}()
	for _, v := range [][2]string{{"@sbx-managed", "1"}, {"@sbx-name", name}, {"@sbx-workspace", workspace}, {"@sbx-docker-name", name}, {"@sbx-role", "agent"}, {"@sbx-last-session", last}} {
		if e := app.tmuxRun(ctx, "set-option", "-t", name, v[0], v[1]); e != nil {
			return e
		}
	}
	if e := app.tmuxRun(ctx, "set-window-option", "-t", "="+name+":agent", "automatic-rename", "off"); e != nil {
		return e
	}
	if e := app.setTmuxWorkspaceStatus(ctx, name, name); e != nil {
		return e
	}
	if e := app.configureTmuxEnvironment(ctx); e != nil {
		return e
	}
	configured = true
	return nil
}
func (app *App) restoreAgentTmuxSession(ctx context.Context, name string) error {
	if e := validateName(name); e != nil {
		return e
	}
	if app.tmuxSessionExists(ctx, name) {
		return fmt.Errorf("refusing existing unmanaged or invalid session: %s", name)
	}
	path, e := app.assertSafeWorkspace(filepath.Join(app.config.WorkspaceRoot, name), name)
	if e != nil {
		return e
	}
	exists, e := app.nativeSandboxNamedExists(ctx, name, false)
	if e != nil {
		return e
	}
	if !exists {
		return fmt.Errorf("Docker Sandbox %s is missing; preserving workspace", name)
	}
	last := name
	for _, s := range app.listTmuxSessions(ctx) {
		if s.Managed && s.WorkspaceName == name && s.Workspace == path && s.Role == "project" {
			last = s.Name
			break
		}
	}
	return app.createAgentTmuxSession(ctx, name, path, last, "")
}
func (app *App) openSession(ctx context.Context, name string) error {
	if e := validateName(name); e != nil {
		return e
	}
	if _, e := app.assertSafeWorkspace(filepath.Join(app.config.WorkspaceRoot, name), name); e != nil {
		return e
	}
	if !app.tmuxSessionManaged(ctx, name) {
		if e := app.restoreAgentTmuxSession(ctx, name); e != nil {
			return e
		}
	}
	target := app.tmuxOption(ctx, name, "@sbx-last-session")
	if target == "" || !app.tmuxSessionExists(ctx, target) || app.tmuxOption(ctx, target, "@sbx-name") != name {
		target = name
	}
	return app.attachOrSwitch(ctx, target)
}
func (app *App) configureTmuxEnvironment(ctx context.Context) error {
	for _, v := range [][2]string{{"SBX_WORKFLOW_BIN", app.Self}, {"SBX_CONFIG_DIR", app.config.ConfigDir}, {"SBX_CONFIG_FILE", app.config.ConfigFile}, {"SBX_HOST_SHELL", app.config.HostShell}} {
		if e := app.tmuxRun(ctx, "set-environment", "-g", v[0], v[1]); e != nil {
			return e
		}
	}
	for _, v := range [][2]string{{"default-shell", "/bin/bash"}, {"default-command", "exec " + process.ShellQuote(app.config.HostShell) + " -l"}, {"remain-on-exit", "failed"}, {"status-left-length", "50"}} {
		if e := app.tmuxRun(ctx, "set-option", "-g", v[0], v[1]); e != nil {
			return e
		}
	}
	for _, b := range [][2]string{{"i", "new"}, {"o", "tui"}, {"f", "tui --current"}, {"s", "tui --current --split"}, {"a", "add"}} {
		// run-shell expands formats at keypress time. Popup environment values and
		// its shell command otherwise receive these format expressions literally.
		popup := append([]string{"tmux"}, app.tmuxArgs("display-popup", "-E", "-w", "90%", "-h", "90%")...)
		for i := range popup {
			popup[i] = process.ShellQuote(popup[i])
		}
		command := strings.Join(popup, " ") +
			` -c #{q:client_name} -t #{q:pane_id}` +
			` -e SBX_TMUX_CLIENT=#{q:client_name} -e SBX_TMUX_PANE=#{q:pane_id} ` +
			process.ShellQuote(`exec "$SBX_WORKFLOW_BIN" __popup `+b[1])
		if e := app.tmuxRun(ctx, "bind-key", b[0], "run-shell", "-b", command); e != nil {
			return e
		}
	}
	if e := app.tmuxRun(ctx, "bind-key", "n", "select-pane", "-t", ".+"); e != nil {
		return e
	}
	reload, _ := app.tmuxOutput(ctx, true, "list-keys", "-T", "prefix", "r")
	if strings.Contains(reload, "source-file") || strings.Contains(reload, "tmux-reload") {
		if e := app.tmuxRun(ctx, "bind-key", "r", "run-shell", `"$SBX_WORKFLOW_BIN" tmux-reload`); e != nil {
			return e
		}
	}
	for _, key := range []string{"C-h", "C-l"} {
		old, _ := app.tmuxOutput(ctx, true, "list-keys", "-T", "root", key)
		if old == "" {
			dir := "-L"
			if key == "C-l" {
				dir = "-R"
			}
			if e := app.tmuxRun(ctx, "bind-key", "-n", key, "select-pane", dir); e != nil {
				return e
			}
		}
	}
	return app.configureSidecarNumberBindings(ctx)
}
func (app *App) cmdTmuxReload(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Usage: sbx tmux-reload")
	}
	path := app.config.TmuxConfig
	if path == "" {
		h, e := os.UserHomeDir()
		if e != nil {
			return e
		}
		path = filepath.Join(h, ".tmux.conf")
	}
	if _, e := os.Stat(path); e == nil {
		if e = app.tmuxRun(context.Background(), "source-file", path); e != nil {
			return e
		}
	}
	return app.configureTmuxEnvironment(context.Background())
}
func nestedTmuxCommand(args []string) string {
	quoted := make([]string, len(args))
	for i, s := range args {
		quoted[i] = process.ShellQuote(s)
	}
	return "exec env -u TMUX -u TMUX_PANE tmux " + strings.Join(quoted, " ")
}
func prefixBindingCommand(output, key string) (string, bool) {
	fields := strings.Fields(output)
	for i := 0; i+2 < len(fields); i++ {
		if fields[i] == "-T" && fields[i+1] == "prefix" && fields[i+2] == key {
			prefix := strings.Join(fields[:i+3], " ")
			_ = prefix
			needle := " " + key + " "
			p := strings.Index(output, needle)
			if p >= 0 {
				return strings.TrimSpace(output[p+len(needle):]), true
			}
		}
	}
	return "", false
}
func (app *App) configureSidecarNumberBindings(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		k := strconv.Itoa(i)
		v, _ := app.tmuxOutput(ctx, true, "list-keys", "-T", "prefix", k)
		if strings.Contains(v, "__sidecar-key") {
			continue
		}
		old, ok := prefixBindingCommand(v, k)
		if !ok {
			continue
		}
		if e := app.tmuxRun(ctx, "set-option", "-g", "@sbx-key-"+k, old); e != nil {
			return e
		}
		wrapper := `"$SBX_WORKFLOW_BIN" __sidecar-key ` + k + ` '#{pane_id}'`
		if e := app.tmuxRun(ctx, "bind-key", k, "if-shell", "-F", "#{@sbx-sidecar}", "run-shell "+process.ShellQuote(wrapper), old); e != nil {
			return e
		}
	}
	return nil
}
func (app *App) cmdSidecarKey(args []string) error {
	if len(args) != 2 || len(args[0]) != 1 || args[0][0] < '0' || args[0][0] > '9' || !strings.HasPrefix(args[1], "%") {
		return fmt.Errorf("invalid sidecar key")
	}
	ctx := context.Background()
	side, e := app.tmuxOutput(ctx, true, "show-options", "-p", "-qv", "-t", args[1], "@sbx-sidecar")
	if e != nil || trimCommandOutput(side) != "1" {
		return fmt.Errorf("not a sidecar pane")
	}
	tty, e := app.tmuxOutput(ctx, true, "display-message", "-p", "-t", args[1], "#{pane_tty}")
	if e != nil {
		return e
	}
	out, e := app.tmuxOutput(ctx, true, "list-clients", "-F", "#{client_tty}\t#{session_name}")
	if e != nil {
		return e
	}
	for _, l := range outputLines(out) {
		f := strings.Split(l, "\t")
		if len(f) == 2 && f[0] == trimCommandOutput(tty) {
			return app.tmuxRun(ctx, "select-window", "-t", "="+f[1]+":"+args[0])
		}
	}
	return fmt.Errorf("sidecar client not attached")
}
func (app *App) cmdSplit(args []string) error {
	if len(args) == 0 {
		return app.cmdTUI([]string{"--current", "--split"})
	}
	if len(args) > 1 {
		return fmt.Errorf("Usage: sbx split [project|agent]")
	}
	ctx := context.Background()
	c, e := app.currentManagedContext(ctx)
	if e != nil {
		return e
	}
	sel := model.Selection{Kind: model.SelectProject, Workspace: c.WorkspaceName, Project: args[0]}
	if args[0] == "agent" {
		sel.Kind = model.SelectAgent
		sel.Project = ""
	}
	return app.splitSelection(ctx, c, sel)
}
func (app *App) splitSelection(ctx context.Context, c managedContext, s model.Selection) error {
	if s.Workspace != c.WorkspaceName {
		return fmt.Errorf("split must belong to current workspace")
	}
	target := s.Workspace
	if s.Kind == model.SelectProject {
		projects, e := app.loadProjects()
		if e != nil {
			return e
		}
		if !app.projectConfigured(s.Project, projects) {
			return fmt.Errorf("unknown project")
		}
		real, e := canonicalDirectory(filepath.Join(c.Workspace, s.Project))
		if e != nil || filepath.Dir(real) != c.Workspace {
			return fmt.Errorf("project is missing or outside workspace")
		}
		target, e = app.ensureProjectSession(ctx, c, s.Project, false)
		if e != nil {
			return e
		}
	} else if s.Kind != model.SelectAgent {
		return fmt.Errorf("invalid split selection")
	}
	if target == c.Session {
		return fmt.Errorf("cannot split session into itself")
	}
	if !app.tmuxSessionExists(ctx, target) {
		if e := app.restoreAgentTmuxSession(ctx, s.Workspace); e != nil {
			return e
		}
	}
	pane, e := app.tmuxOutput(ctx, false, "split-window", "-h", "-p", "50", "-P", "-F", "#{pane_id}", "-t", c.Pane, nestedTmuxCommand(app.tmuxArgs("attach-session", "-t", "="+target)))
	if e != nil {
		return e
	}
	pane = trimCommandOutput(pane)
	for _, v := range [][2]string{{"@sbx-sidecar", "1"}, {"@sbx-sidecar-target", target}} {
		if e = app.tmuxRun(ctx, "set-option", "-p", "-t", pane, v[0], v[1]); e != nil {
			return e
		}
	}
	return nil
}
