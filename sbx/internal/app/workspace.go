package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sbx/internal/config"
	"sbx/internal/process"
	"strings"
)

func validateName(name string) error {
	if name == "default" {
		return errors.New("'default' is reserved by Docker Sandboxes")
	}
	if name == "" || !sandboxNamePattern.MatchString(name) {
		return fmt.Errorf("invalid sandbox name '%s' (use at least two letters, numbers, dots, or hyphens; start with a letter or number)", name)
	}
	if len(name) < 2 {
		return errors.New("sandbox name must contain at least two characters")
	}
	return nil
}

func (app *App) promptLine(prompt string) (string, error) {
	fmt.Fprint(app.Err, prompt)
	reader := bufio.NewReader(app.In)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	if err != nil {
		return "", errors.New("no name provided")
	}
	value = strings.TrimSuffix(value, "\n")
	value = strings.TrimSuffix(value, "\r")
	if value == "" {
		return "", errors.New("name must not be empty")
	}
	return value, nil
}

func (app *App) promptValue(prompt string) (string, error) {
	fmt.Fprint(app.Err, prompt)
	reader, ok := app.In.(*bufio.Reader)
	if !ok {
		reader = bufio.NewReader(app.In)
	}
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("cannot read workspace input: %w", err)
	}
	return strings.TrimSuffix(strings.TrimSuffix(value, "\n"), "\r"), nil
}

func (app *App) loadProjects() ([]config.Project, error) {
	return config.LoadProjects(app.config.ProjectsFile)
}

func (app *App) projectConfigured(name string, projects []config.Project) bool {
	for _, project := range projects {
		if project.Name == name {
			return true
		}
	}
	return false
}

func (app *App) chooseWithFZF(ctx context.Context, prompt string, items []string, multi bool) ([]string, bool, error) {
	if err := app.checkCommand("fzf"); err != nil {
		return nil, false, err
	}
	args := make([]string, 0, 8)
	if multi {
		args = append(args, "--multi")
	}
	args = append(args, "--prompt", prompt)
	if multi {
		args = append(args, "--header", "Space select, Enter confirm", "--bind", "space:toggle,tab:down,btab:up")
	}
	args = append(args, "--height=100%", "--layout=reverse", "--border")
	input := ""
	if len(items) > 0 {
		input = strings.Join(items, "\n") + "\n"
	}
	output, err := app.runner.Output(ctx, "fzf", args, process.Options{Stdin: strings.NewReader(input)})
	if err != nil {
		return nil, false, nil
	}
	selection := outputLines(output)
	if len(selection) == 0 {
		return nil, false, nil
	}
	return selection, true, nil
}

func (app *App) assertSafeWorkspace(workspace, name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}
	rootReal, err := canonicalDirectory(app.config.WorkspaceRoot)
	if err != nil {
		return "", fmt.Errorf("workspace root does not exist: %s", app.config.WorkspaceRoot)
	}
	workspaceReal, err := canonicalDirectory(workspace)
	if err != nil {
		return "", fmt.Errorf("workspace does not exist: %s", workspace)
	}
	expected := filepath.Join(rootReal, name)
	if workspaceReal != expected {
		return "", fmt.Errorf("refusing unsafe workspace path: %s", workspaceReal)
	}
	if !markerMatches(workspaceReal, name) {
		return "", fmt.Errorf("refusing unmarked workspace: %s", workspaceReal)
	}
	return workspaceReal, nil
}

func (app *App) cloneProjectFromCache(ctx context.Context, project config.Project, destination string) error {
	if err := app.updateProjectGitCache(ctx, project); err != nil {
		return fmt.Errorf("Git cache update failed: %s", project.Name)
	}
	cache := filepath.Join(app.config.GitCacheRoot, project.Name+".git")
	app.note("- cloning %s from cache", project.Name)
	if err := app.runner.Run(ctx, "git", []string{"clone", "--branch", project.Branch, "--", cache, destination}, process.Options{}); err != nil {
		return fmt.Errorf("clone failed: %s", project.Name)
	}
	if err := app.runner.Run(ctx, "git", []string{"-C", destination, "remote", "set-url", "origin", project.Repo}, process.Options{}); err != nil {
		return fmt.Errorf("cannot set origin for project: %s", project.Name)
	}
	return nil
}

func (app *App) bootstrapProject(ctx context.Context, project config.Project, path string) error {
	useNode, err := projectDeclaresNodeVersion(path)
	if err != nil {
		return fmt.Errorf("Node version configuration for %s: %w", project.Name, err)
	}
	options := process.Options{Dir: path}
	if useNode {
		if err := app.checkCommand("fnm"); err != nil {
			return err
		}
		// Resolve only this repository's settings, never a parent workspace pin.
		options.Env = map[string]string{"FNM_VERSION_FILE_STRATEGY": "local", "FNM_RESOLVE_ENGINES": "true"}
		app.note("- installing repository Node version with fnm for %s", project.Name)
		if err := app.runner.Run(ctx, "fnm", []string{"install"}, options); err != nil {
			return fmt.Errorf("Node version setup failed: %s: %w", project.Name, err)
		}
	}
	runNode := func(command string, args []string) error {
		if useNode {
			return app.runner.Run(ctx, "fnm", append([]string{"exec", "--", command}, args...), options)
		}
		return app.runner.Run(ctx, command, args, options)
	}
	if project.BeforeBootstrap != "" {
		app.note("- running pre-bootstrap for %s", project.Name)
		if err := runNode(app.config.HostShell, []string{"-c", project.BeforeBootstrap}); err != nil {
			return fmt.Errorf("pre-bootstrap failed: %s", project.Name)
		}
	}
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		if err := app.checkCommand("go"); err != nil {
			return err
		}
		app.note("- preparing Go dependencies for %s", project.Name)
		for _, command := range [][]string{{"mod", "download"}, {"mod", "vendor"}, {"mod", "tidy"}} {
			if err := app.runner.Run(ctx, "go", command, process.Options{Dir: path}); err != nil {
				return fmt.Errorf("go %s failed: %s", strings.Join(command, " "), project.Name)
			}
		}
	}
	if _, err := os.Stat(filepath.Join(path, "pnpm-lock.yaml")); err == nil {
		if !useNode {
			if err := app.checkCommand("pnpm"); err != nil {
				return err
			}
		}
		app.note("- installing pnpm dependencies for %s", project.Name)
		if err := runNode("pnpm", []string{"install"}); err != nil {
			return fmt.Errorf("pnpm install failed: %s", project.Name)
		}
	}
	return nil
}

func projectDeclaresNodeVersion(path string) (bool, error) {
	for _, name := range []string{".node-version", ".nvmrc"} {
		contents, err := os.ReadFile(filepath.Join(path, name))
		if err == nil {
			if strings.TrimSpace(string(contents)) == "" {
				return false, fmt.Errorf("%s is empty", name)
			}
			return true, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
	}
	contents, err := os.ReadFile(filepath.Join(path, "package.json"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var manifest struct {
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`
	}
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return false, err
	}
	return strings.TrimSpace(manifest.Engines.Node) != "", nil
}

var sandboxNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.-]*$`)
