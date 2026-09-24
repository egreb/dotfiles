package app

import (
	"sbx/internal/model"
	"testing"
)

func TestPreviousWorkspace(t *testing.T) {
	inventory := model.Inventory{Workspaces: []model.Workspace{{Name: "deleted", LastUsed: 100}, {Name: "alpha", LastUsed: 10}, {Name: "beta", LastUsed: 20}, {Name: "unknown"}}}
	if got := previousWorkspace(inventory, "deleted"); got != "beta" {
		t.Fatal(got)
	}
	inventory.Workspaces[1].LastUsed = 20
	if got := previousWorkspace(inventory, "deleted"); got != "alpha" {
		t.Fatal("ties must use names", got)
	}
	if got := previousWorkspace(model.Inventory{Workspaces: inventory.Workspaces[:1]}, "deleted"); got != "" {
		t.Fatal("no remaining workspace", got)
	}
}
