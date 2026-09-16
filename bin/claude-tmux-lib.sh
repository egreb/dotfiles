#!/usr/bin/env bash
# Shared helpers for surfacing Claude Code session state in the tmux launchers
# (tmux-sessionizer, tmux-open-session, ...). Source this file; don't run it.
#
# Presence is derived from the live `claude` process (self-healing); the exact
# state (ready / needs feedback) is refined by ~/bin/claude-tmux-status, written
# from Claude Code hooks. Icons are overridable via the CLAUDE_TMUX_ICON_* vars.

CLAUDE_TMUX_STATUS_DIR="${CLAUDE_TMUX_STATUS_DIR:-$HOME/.cache/claude-tmux-status}"
CLAUDE_TMUX_ICON_READY="${CLAUDE_TMUX_ICON_READY:-🟢}"
CLAUDE_TMUX_ICON_FEEDBACK="${CLAUDE_TMUX_ICON_FEEDBACK:-🔴}"
CLAUDE_TMUX_ICON_WORKING="${CLAUDE_TMUX_ICON_WORKING:-🟡}"
# Placeholder shown when a line has no icon, so text stays aligned. Should be as
# wide as an icon (the default emoji render as two cells -> two spaces).
CLAUDE_TMUX_ICON_BLANK="${CLAUDE_TMUX_ICON_BLANK:-  }"

# Emit "present\t<cwd>" for every directory a live `claude` process runs in.
claude_present_records() {
    local pids
    pids=$(ps -axo pid=,comm= 2>/dev/null | awk '
        { pid=$1; $1=""; c=$0; sub(/^[ \t]+/, "", c)
          n=split(c, a, "/"); if (a[n] == "claude") print pid }')
    [[ -z $pids ]] && return 0
    lsof -a -d cwd -w -Fn -p "$(echo "$pids" | paste -sd, -)" 2>/dev/null |
        awk '/^n\// { print "present\t" substr($0, 2) }'
}

# Emit "<state>\t<path>" for each fresh (<24h) hook-written state file.
claude_state_records() {
    find "$CLAUDE_TMUX_STATUS_DIR" -type f -mmin -1440 -exec cat {} + 2>/dev/null
}
