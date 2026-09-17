package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigDir       string
	ConfigFile      string
	WorkspaceRoot   string
	ProjectsFile    string
	GitCacheRoot    string
	Agent           string
	ClaudeTheme     string
	HostShell       string
	TmuxSocket      string
	TmuxConfig      string
	DockerSBXBinary string
}

type Project struct {
	Name            string
	Repo            string
	Branch          string
	BeforeBootstrap string
}

var projectNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func Load(_ string) (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("cannot determine home directory: %w", err)
	}

	configDir := os.Getenv("SBX_CONFIG_DIR")
	if configDir == "" {
		configDir = filepath.Join(home, ".config", "sbx")
	}
	configDir = expandHome(configDir, home)

	configFile := os.Getenv("SBX_CONFIG_FILE")
	if configFile == "" {
		configFile = filepath.Join(configDir, "config.yaml")
	}
	if info, err := os.Stat(configFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("configuration not found: %s (copy config.yaml.example and projects.yaml.example into %s)", configFile, configDir)
		}
		return Config{}, fmt.Errorf("cannot read configuration %s: %w", configFile, err)
	} else if !info.Mode().IsRegular() {
		return Config{}, fmt.Errorf("configuration not found: %s (copy config.yaml.example and projects.yaml.example into %s)", configFile, configDir)
	}

	values, err := loadStringMap(configFile, "configuration")
	if err != nil {
		return Config{}, err
	}

	get := func(key, fallback string) (string, error) {
		node, ok := values[key]
		if !ok {
			return fallback, nil
		}
		if node.Kind != yaml.ScalarNode || (node.Tag != "!!str" && node.Tag != "!!int" && node.Tag != "!!float") {
			return "", fmt.Errorf("configuration value %q must be a string or number", key)
		}
		if node.Tag == "!!str" {
			return node.Value, nil
		}
		// Match Ruby's puts for ordinary YAML numeric scalars.
		if node.Tag == "!!int" {
			value, parseErr := strconv.ParseInt(node.Value, 0, 64)
			if parseErr == nil {
				return strconv.FormatInt(value, 10), nil
			}
		}
		if node.Tag == "!!float" {
			value, parseErr := strconv.ParseFloat(node.Value, 64)
			if parseErr == nil {
				return strconv.FormatFloat(value, 'g', -1, 64), nil
			}
		}
		return node.Value, nil
	}

	workspaceRoot, err := get("workspace_root", filepath.Join(home, "sandboxes"))
	if err != nil {
		return Config{}, err
	}
	projectsFile, err := get("projects_file", filepath.Join(configDir, "projects.yaml"))
	if err != nil {
		return Config{}, err
	}
	gitCacheRoot, err := get("git_cache_root", filepath.Join(home, ".cache", "sbx", "git"))
	if err != nil {
		return Config{}, err
	}
	agent, err := get("agent", "claude")
	if err != nil {
		return Config{}, err
	}
	claudeTheme, err := get("claude_theme", "light")
	if err != nil {
		return Config{}, err
	}
	switch claudeTheme {
	case "", "light", "dark", "light-daltonized", "dark-daltonized", "light-ansi", "dark-ansi":
	default:
		return Config{}, fmt.Errorf("invalid claude_theme: %s", claudeTheme)
	}
	hostShell, err := get("host_shell", "fish")
	if err != nil {
		return Config{}, err
	}
	tmuxSocket, err := get("tmux_socket", "sbx-v1")
	if err != nil {
		return Config{}, err
	}
	tmuxConfig, err := get("tmux_config", "")
	if err != nil {
		return Config{}, err
	}
	dockerBinary, err := get("docker_sbx_bin", "")
	if err != nil {
		return Config{}, err
	}

	workspaceRoot = expandHome(workspaceRoot, home)
	projectsFile = expandHome(projectsFile, home)
	gitCacheRoot = expandHome(gitCacheRoot, home)
	tmuxConfig = expandHome(tmuxConfig, home)
	dockerBinary = expandHome(dockerBinary, home)
	hostShell = expandHome(hostShell, home)

	if strings.ContainsRune(hostShell, filepath.Separator) {
		info, statErr := os.Stat(hostShell)
		if statErr != nil || info.IsDir() || info.Mode()&0o111 == 0 {
			return Config{}, fmt.Errorf("host_shell is not executable: %s", hostShell)
		}
	} else {
		resolved, lookupErr := exec.LookPath(hostShell)
		if lookupErr != nil {
			return Config{}, errors.New("host_shell command not found")
		}
		hostShell = resolved
	}

	if agent != "claude" && agent != "codex" {
		return Config{}, fmt.Errorf("agent must be 'claude' or 'codex', got: %s", agent)
	}
	if !validTmuxSocket(tmuxSocket) {
		return Config{}, fmt.Errorf("invalid tmux_socket: %s", tmuxSocket)
	}
	if tmuxConfig != "" {
		if info, statErr := os.Stat(tmuxConfig); statErr != nil || !info.Mode().IsRegular() {
			return Config{}, fmt.Errorf("tmux configuration not found: %s", tmuxConfig)
		}
	}

	return Config{
		ConfigDir:       configDir,
		ConfigFile:      configFile,
		WorkspaceRoot:   workspaceRoot,
		ProjectsFile:    projectsFile,
		GitCacheRoot:    gitCacheRoot,
		Agent:           agent,
		ClaudeTheme:     claudeTheme,
		HostShell:       hostShell,
		TmuxSocket:      tmuxSocket,
		TmuxConfig:      tmuxConfig,
		DockerSBXBinary: dockerBinary,
	}, nil
}
