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

pane_dir=$(tmux display-message -p '#{pane_current_path}')
pane_id=$(tmux display-message -p '#{pane_id}')
pane_pid=$(tmux display-message -p '#{pane_pid}')

# Detect AI tools running in pane
ai_mode=false
pgrep -P "$pane_pid" -f "opencode|claude|codex" >/dev/null && ai_mode=true

# Find git root (fallback to pane_dir)
git_root=$(cd "$pane_dir" && git rev-parse --show-toplevel 2>/dev/null || echo "$pane_dir")

# Pick files with fd + fzf + bat preview
selected=$(
  cd "$git_root" && fd --type f --hidden --follow --exclude .git | \
    fzf --multi --height 100% --preview "bat --theme='Catppuccin Mocha' --style=numbers --color=always {}"
)

[ -z "$selected" ] && exit 0

# Format and send to pane
if $ai_mode; then
  printf -v output "@%s " $(echo "$selected")
else
  output=$(printf "%q " $(echo "$selected"))
fi

tmux send-keys -t "$pane_id" "$output"
