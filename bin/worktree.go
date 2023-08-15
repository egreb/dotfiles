///usr/bin/env go run "$0" "$@"; exit

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const project = "~/code/work/cybr/baloo-frontend-platform/master"

func main() {
	cmd := exec.Command("git", "-C", project, "worktree", "list")
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("failed to run git command", err)
		os.Exit(1)
	}

	list := strings.Split(string(b), "\n\n")

	for _, l := range list {
		if len(l) == 0 {
			continue
		}

		data := strings.Split(l, "\n")
		// path := strings.Split(data[0], " ")[1]
		// hash := strings.Split(data[1], " ")[1]
		branch := strings.Split(data[2], " ")[1]

		mergeBaseCmd := exec.Command("git", "-C", project, "merge-base", "master", branch)
		revParseCmd := exec.Command("git", "-C", project, "rev-parse", branch)

		mergeBaseRes, err := mergeBaseCmd.Output()
		revParseRes, err := revParseCmd.Output()
		if err != nil {
			log.Fatalln("failed to execute git commands", err)
		}

		if string(mergeBaseRes) == string(revParseRes) {
			fmt.Printf("branch %s is merged\n", branch)
		} else {
			fmt.Printf("branch %s is not merged\n", branch)
		}
	}

	os.Exit(1)
}
