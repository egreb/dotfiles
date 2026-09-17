package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type Runner struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Options struct {
	Dir      string
	Env      map[string]string
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	QuietErr bool
}

func (runner Runner) Run(ctx context.Context, name string, args []string, options Options) error {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = options.Dir
	command.Env = mergeEnvironment(os.Environ(), options.Env)
	command.Stdin = firstReader(options.Stdin, runner.In)
	command.Stdout = firstWriter(options.Stdout, runner.Out)
	if options.QuietErr {
		command.Stderr = io.Discard
	} else {
		command.Stderr = firstWriter(options.Stderr, runner.Err)
	}
	return command.Run()
}

func (runner Runner) Output(ctx context.Context, name string, args []string, options Options) (string, error) {
	var stdout bytes.Buffer
	options.Stdout = &stdout
	err := runner.Run(ctx, name, args, options)
	return stdout.String(), err
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}
	return 1
}

func Replace(binary string, args []string) error {
	return syscall.Exec(binary, append([]string{binary}, args...), os.Environ())
}

func CanonicalExecutable(target string) (string, error) {
	if !filepath.IsAbs(target) && !strings.ContainsRune(target, filepath.Separator) {
		resolved, err := exec.LookPath(target)
		if err != nil {
			return "", err
		}
		target = resolved
	}
	if !filepath.IsAbs(target) {
		absolute, err := filepath.Abs(target)
		if err != nil {
			return "", err
		}
		target = absolute
	}
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		return "", err
	}
	return filepath.Abs(resolved)
}

func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}

func FindNextOnPath(name, self string) (string, error) {
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		candidate := filepath.Join(dir, name)
		if !IsExecutable(candidate) {
			continue
		}
		real, err := CanonicalExecutable(candidate)
		if err == nil && real != self {
			return real, nil
		}
	}
	return "", fmt.Errorf("executable not found: %s", name)
}
func ShellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
func firstReader(a, b io.Reader) io.Reader {
	if a != nil {
		return a
	}
	return b
}
func firstWriter(a, b io.Writer) io.Writer {
	if a != nil {
		return a
	}
	return b
}
func mergeEnvironment(base []string, override map[string]string) []string {
	result := make([]string, 0, len(base)+len(override))
	for _, item := range base {
		k, _, _ := strings.Cut(item, "=")
		if _, ok := override[k]; !ok {
			result = append(result, item)
		}
	}
	for k, v := range override {
		result = append(result, k+"="+v)
	}
	return result
}
