#!/bin/bash

# Fuzzy search git worktrees by name or branch
function ot() {
	selected=$(git worktree list | fzf --border --no-header --margin='5,5' --layout=reverse-list)
	dir=$(echo $selected | awk '{split($0,a); print $1}')
	branch=$(echo $selected | awk '{split($0,a); print $3}')

	if ! [ -z "$dir" ]; then
		tmux new-window -n $branch -c $dir
	else
		echo "Not found"
	fi
}
ot "$@"
