package model

type Session struct {
	Name, Role, Project, ActivePaneID, CurrentCommand string
	Attached, PaneDead                                bool
}
type Project struct {
	Name, Path string
	Session    *Session
}
type Workspace struct {
	Name, Path   string
	LastUsed     int64
	AgentSession *Session
	Projects     []Project
}
type Inventory struct{ Workspaces []Workspace }
type SelectionKind int

const (
	SelectNone SelectionKind = iota
	SelectWorkspace
	SelectAgent
	SelectProject
)

type Selection struct {
	Kind               SelectionKind
	Workspace, Project string
}
