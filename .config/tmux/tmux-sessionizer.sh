#!/usr/bin/env bash
set -euo pipefail

# Prefer the installed config, falling back to the repo before Stow is rerun.
export FZF_DEFAULT_OPTS_FILE=""
for fzf_config in "$HOME/.config/fzf/config" "$HOME/dotfiles/.config/fzf/config"; do
  if [[ -r "$fzf_config" ]]; then
    export FZF_DEFAULT_OPTS_FILE="$fzf_config"
    break
  fi
done

paths="${TMUX_SESSIONIZER_PATHS:-$HOME}"
default_depth="${TMUX_SESSIONIZER_DEPTH:-1}"

expand_tilde() { echo "${1/#~/$HOME}"; }

selected=$(
  {
    # Existing tmux sessions
    tmux list-sessions -F '[TMUX] #{session_name}' 2>/dev/null || true

    # Discovered directories from paths
    for entry in $paths; do
      # Extract optional depth suffix (e.g., ~/foo:2)
      [[ "$entry" =~ ^([^:]+):([0-9]+)$ ]] && path="${BASH_REMATCH[1]}" depth="${BASH_REMATCH[2]}" || { path="$entry"; depth="$default_depth"; }
      path=$(expand_tilde "$path")
      for expanded in $path; do
        [ -d "$expanded" ] || continue
        find "$expanded" -mindepth 1 -maxdepth "$depth" -type d | sed "s|^$HOME|~|"
      done
    done
  } | fzf --height 100%
)

[ -z "$selected" ] && exit 0
selected=$(expand_tilde "$selected")

# If existing session selected, switch to it
if [[ "$selected" =~ ^\[TMUX\]\ (.+)$ ]]; then
  sess="${BASH_REMATCH[1]}"
  [ -z "${TMUX:-}" ] && tmux attach -t "$sess" || tmux switch-client -t "$sess"
  exit 0
fi

# Otherwise, create and switch/attach
sess=$(basename "$selected")
if [ -z "${TMUX:-}" ]; then
  tmux has-session -t "$sess" 2>/dev/null || tmux new-session -ds "$sess" -c "$selected"
  tmux attach -t "$sess"
else
  tmux has-session -t "$sess" 2>/dev/null || tmux new-session -ds "$sess" -c "$selected"
  tmux switch-client -t "$sess"
fi
