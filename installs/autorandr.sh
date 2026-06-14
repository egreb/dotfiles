set -e

# autorandr remembers a display layout per monitor combination and reapplies it
# automatically on hotplug (it installs udev rules + systemd user units). i3
# also calls `autorandr --change` on login/restart (see .config/i3/config).
#
# This solves multi-monitor docking: workspaces, positions, and which output is
# primary (and thus where the i3 tray lives) all follow the saved profile, so
# undocking while logged in no longer strands the tray on a dead output.

echo "installing autorandr"
sudo apt install -y autorandr

cat <<'EOF'

autorandr installed. Now save one profile per setup (run these once, in i3):

  # Docked: external as the main/primary screen, right of the laptop.
  xrandr --output DP-1 --primary --auto --right-of eDP-1
  autorandr --save docked

  # Laptop only: unplug the external first, then:
  autorandr --save mobile

Check output names with `xrandr --query`. List/apply profiles with
`autorandr` / `autorandr --load <name>`; `autorandr --change` auto-picks by
what's connected. Profiles live in ~/.config/autorandr/ (host-specific, so
they're intentionally NOT tracked in this repo).
EOF
