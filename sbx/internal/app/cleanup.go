package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"sbx/internal/process"
)

// Stop host-side jobs before closing their terminals, while their ancestry still
// identifies the workspace. Never select processes merely by port or command name.
func (app *App) stopWorkspaceJobs(ctx context.Context, sessions []tmuxSession) error {
	roots := map[int]bool{}
	for _, session := range sessions {
		out, err := app.tmuxOutput(ctx, false, "list-panes", "-s", "-t", "="+session.Name, "-F", "#{pane_pid}")
		if err != nil {
			return err
		}
		for _, line := range strings.Fields(out) {
			pid, err := strconv.Atoi(line)
			if err != nil || pid <= 1 {
				return fmt.Errorf("invalid tmux pane PID: %q", line)
			}
			roots[pid] = true
		}
	}
	if len(roots) == 0 {
		return nil
	}
	out, err := app.runner.Output(ctx, "ps", []string{"-axo", "pid=,ppid="}, process.Options{})
	if err != nil {
		return fmt.Errorf("list workspace processes: %w", err)
	}
	parents := map[int]int{}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		pid, e1 := strconv.Atoi(fields[0])
		parent, e2 := strconv.Atoi(fields[1])
		if e1 != nil || e2 != nil {
			return fmt.Errorf("invalid process listing: %q", line)
		}
		parents[pid] = parent
	}
	// Deletion may itself have been invoked from one of these terminals.
	protected := map[int]bool{}
	for pid := os.Getpid(); pid > 1 && !protected[pid]; pid = parents[pid] {
		protected[pid] = true
	}
	jobs := []int{}
	for pid := range parents {
		if pid <= 1 || protected[pid] {
			continue
		}
		if roots[pid] {
			jobs = append(jobs, pid)
			continue
		}
		seen := map[int]bool{}
		for parent := parents[pid]; parent > 1 && !seen[parent]; parent = parents[parent] {
			if roots[parent] {
				jobs = append(jobs, pid)
				break
			}
			seen[parent] = true
		}
	}
	signal := func(pid int, sig syscall.Signal) error {
		err := syscall.Kill(pid, sig)
		if err == syscall.ESRCH {
			return nil
		}
		return err
	}
	for _, pid := range jobs {
		if err := signal(pid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("stop workspace process %d: %w", pid, err)
		}
	}
	deadline := time.Now().Add(2 * time.Second)
	for len(jobs) > 0 {
		alive := jobs[:0]
		for _, pid := range jobs {
			if syscall.Kill(pid, 0) != syscall.ESRCH {
				alive = append(alive, pid)
			}
		}
		jobs = alive
		if len(jobs) == 0 {
			break
		}
		if time.Now().After(deadline) {
			for _, pid := range jobs {
				if err := signal(pid, syscall.SIGKILL); err != nil {
					return fmt.Errorf("kill workspace process %d: %w", pid, err)
				}
			}
			// Signals are asynchronous; wait for exit before ports are reused.
			for until := time.Now().Add(2 * time.Second); ; {
				out, err := app.runner.Output(ctx, "ps", []string{"-axo", "pid=,stat="}, process.Options{})
				if err != nil {
					return err
				}
				live := map[int]bool{}
				for _, line := range strings.Split(out, "\n") {
					fields := strings.Fields(line)
					if len(fields) != 2 {
						continue
					}
					pid, _ := strconv.Atoi(fields[0])
					if !strings.HasPrefix(fields[1], "Z") {
						live[pid] = true
					}
				}
				remaining := false
				for _, pid := range jobs {
					remaining = remaining || live[pid]
				}
				if !remaining {
					return nil
				}
				if time.Now().After(until) {
					return fmt.Errorf("workspace processes have not exited after SIGKILL")
				}
				time.Sleep(50 * time.Millisecond)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
