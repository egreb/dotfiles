set -e

# Make Caps Lock act as Ctrl at the X SERVER level, so it applies to every
# keyboard X sees — including devices that re-enumerate on resume from suspend
# or when you hotplug a keyboard. Doing it here (rather than `setxkbmap` from
# i3) is what makes the remap survive suspend: `setxkbmap` only touches devices
# present when it runs, and i3 doesn't re-run it on resume.
#
# localectl writes /etc/X11/xorg.conf.d/00-keyboard.conf. Args are
# LAYOUT MODEL VARIANT OPTIONS (matching `setxkbmap -query`: us / pc105 / none).

echo "setting X11 keymap: us / pc105 / caps:ctrl_modifier"
sudo localectl set-x11-keymap us pc105 "" caps:ctrl_modifier

echo "written. Verify with: cat /etc/X11/xorg.conf.d/00-keyboard.conf"
echo "Takes effect on next X login (or: sudo systemctl restart display-manager)."
