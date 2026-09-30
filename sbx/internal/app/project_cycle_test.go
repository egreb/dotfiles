package app

import (
	"sbx/internal/model"
	"testing"
)

func TestCycleProject(t *testing.T) {
	projects := []model.Project{{Name: "api"}, {Name: "web"}}
	for _, tc := range []struct {
		session  string
		previous bool
		want     string
	}{
		{"alpha", false, "api"}, {"alpha", true, "web"},
		{"alpha--api", false, "web"}, {"alpha--api", true, "web"},
		{"alpha--web", false, "api"}, {"alpha--web", true, "api"},
		{"other--api", false, "api"},
	} {
		if got := cycleProject(projects, "alpha", tc.session, tc.previous); got != tc.want {
			t.Fatalf("%+v: got %s", tc, got)
		}
	}
	if cycleProject(nil, "alpha", "alpha", false) != "" {
		t.Fatal("empty projects")
	}
	if cycleProject(projects[:1], "alpha", "alpha--api", true) != "api" {
		t.Fatal("single project")
	}
}
