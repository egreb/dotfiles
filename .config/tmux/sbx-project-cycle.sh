#!/usr/bin/env bash
set -euo pipefail

direction="${1:---next}"

sbx project "$direction"

session="$(tmux display-message -p '#S')"
label="${session/--/ › }"

width=$(( ${#label} + 4 > 16 ? ${#label} + 4 : 16 ))

tmux display-popup -E -x 9999 -y 9999 -w "$width" -h 3 -T sbx \
  "printf '%s\n' $(printf '%q' "$label"); sleep 1"
