package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func expandHome(path, home string) string {
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}
func validTmuxSocket(value string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`).MatchString(value)
}
func readYAML(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err = yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	var check func(*yaml.Node) error
	check = func(n *yaml.Node) error {
		if n.Kind == yaml.AliasNode {
			return fmt.Errorf("YAML aliases are not allowed")
		}
		for _, c := range n.Content {
			if err := check(c); err != nil {
				return err
			}
		}
		return nil
	}
	if err = check(&doc); err != nil {
		return nil, err
	}
	if len(doc.Content) != 1 {
		return nil, fmt.Errorf("expected YAML mapping")
	}
	return doc.Content[0], nil
}
func stringMap(n *yaml.Node) (map[string]*yaml.Node, error) {
	if n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected YAML mapping")
	}
	m := map[string]*yaml.Node{}
	for i := 0; i < len(n.Content); i += 2 {
		k := n.Content[i]
		if k.Tag != "!!str" {
			return nil, fmt.Errorf("mapping keys must be strings")
		}
		if _, ok := m[k.Value]; ok {
			return nil, fmt.Errorf("duplicate key %q", k.Value)
		}
		m[k.Value] = n.Content[i+1]
	}
	return m, nil
}
func loadStringMap(path, kind string) (map[string]*yaml.Node, error) {
	n, err := readYAML(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", kind, err)
	}
	return stringMap(n)
}
func LoadProjects(path string) ([]Project, error) {
	m, err := loadStringMap(path, "projects")
	if err != nil {
		return nil, err
	}
	n := m["projects"]
	if n == nil || n.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("projects must be a sequence")
	}
	var result []Project
	seen := map[string]bool{}
	for _, entry := range n.Content {
		fields, err := stringMap(entry)
		if err != nil {
			return nil, err
		}
		get := func(k string) string {
			if v := fields[k]; v != nil && v.Tag == "!!str" {
				return v.Value
			}
			return ""
		}
		p := Project{Name: get("name"), Repo: get("repo"), Branch: get("branch"), BeforeBootstrap: get("before_bootstrap")}
		if !projectNamePattern.MatchString(p.Name) || p.Name == "." || p.Name == ".." || seen[p.Name] {
			return nil, fmt.Errorf("invalid or duplicate project name: %s", p.Name)
		}
		if p.Repo == "" {
			return nil, fmt.Errorf("project %s has no repo", p.Name)
		}
		if p.Branch == "" {
			p.Branch = "main"
		}
		seen[p.Name] = true
		result = append(result, p)
	}
	return result, nil
}
