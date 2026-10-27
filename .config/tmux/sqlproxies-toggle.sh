#!/usr/bin/env bash
set -euo pipefail

# Toggle a dedicated "sqlproxies" window in the current session. The window
# keeps running when you toggle away, so the proxies stay up.
name="sqlproxies"

if [[ "$(tmux display-message -p '#{window_name}')" == "$name" ]]; then
  tmux last-window
elif tmux list-windows -F '#{window_name}' | grep -qx "$name"; then
  tmux select-window -t ":=$name"
else
  # Go through fish so PATH matches an interactive shell.
  tmux new-window -n "$name" "fish -c 'exec sqlproxies'"
fi
