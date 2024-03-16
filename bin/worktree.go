///usr/bin/env go run "$0" "$@"; exit

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type git struct {
	project string
}

func (g *git) getWorkTreeListRaw() ([]string, error) {
	cmd := exec.Command("git", "-C", g.project, "worktree", "list")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git.GetWorktrees failed: %w", err)
	}
	list := strings.Split(string(b), "\n\n")
	list = strings.Split(string(b), "\n")

	return list, nil
}

type Worktree struct {
	path   string
	branch string
}

func (g *git) getWorkTreeList() ([]Worktree, error) {
	list, err := g.getWorkTreeListRaw()
	if err != nil {
		return nil, fmt.Errorf("g.parseWorkTree failed: %w", err)
	}

	var ww []Worktree
	for _, l := range list {
		data := strings.Split(l, " ")
		path := data[0]
		branch := data[len(data)-1]
		branch = strings.TrimLeft(strings.TrimRight(branch, "]"), "[")

		ww = append(ww, Worktree{path, branch})
	}

	return ww, nil
}

const yes = "y"

func userSaidYes(input string) bool {
	if input == "" {
		return false
	}

	inputLowercase := strings.ToLower(input)
	return inputLowercase == yes
}

func main() {
	project, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to get current working dir", err)
	}
	g := git{project}

	list, err := g.getWorkTreeList()
	if err != nil {
		log.Fatalln("failed to get worktrees", err)
	}

	for _, l := range list {
		if l.branch == "master" || l.branch == "main" {
			continue
		}
		mergeBaseCmd := exec.Command("git", "-C", project, "merge-base", "master", l.branch)
		revParseCmd := exec.Command("git", "-C", project, "rev-parse", l.branch)

		mergeBaseRes, err := mergeBaseCmd.Output()
		if err != nil {
			log.Println("failed to execute git merge base command", err)
			continue
		}
		revParseRes, err := revParseCmd.Output()
		if err != nil {
			log.Println("failed to execute git rev parse command", err)
			continue
		}

		if string(mergeBaseRes) == string(revParseRes) {
			fmt.Printf("branch %s is merged\n", l.branch)
			fmt.Println("deleting worktree", l.path)
			rmWorktreeCmd := exec.Command("git", "-C", project, "worktree", "remove", l.path)
			_, err := rmWorktreeCmd.Output()
			if err != nil {
				log.Println("failed to remove worktree", l.path)
			} else {
				log.Println("removed worktree", l.path)
			}

		} else {
			fmt.Printf("branch %s is not merged\n", l.branch)
			var input string
			fmt.Println("Delete anyway?:")
			_, err := fmt.Scanln(&input)
			if err != nil {
				log.Println("soemthing went wrong asking for user input:", err)
				continue
			}
			if userSaidYes(input) {
				rmWorktreeCmd := exec.Command("git", "-C", project, "worktree", "remove", "--force", l.path)
				_, err := rmWorktreeCmd.Output()
				if err != nil {
					log.Println("failed to remove worktree", l.path)
				} else {
					log.Println("removed worktree", l.path)
				}
			}
		}
	}

	os.Exit(1)
}
