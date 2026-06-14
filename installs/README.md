# installs

Per-tool setup scripts for a fresh machine. Run them individually, e.g.:

```bash
./installs/brew.sh
./installs/i3.sh
```

## Gotchas

### Caps Lock → Ctrl remap lost after suspend

If the remap lives in i3 as `exec_always setxkbmap -option caps:ctrl_modifier`,
it can silently revert after resuming from suspend (or after plugging in a
keyboard). `setxkbmap` only configures the keyboards present when it runs, and
on resume USB keyboards often re-enumerate — X then applies its *default* XKB
config to the re-added device, and i3 never re-runs `exec_always` on resume.

Fix (handled by `keyboard.sh`): set the option at the X server level so it
applies to every keyboard, including ones that re-appear:

```bash
sudo localectl set-x11-keymap us pc105 "" caps:ctrl_modifier
```

This writes `/etc/X11/xorg.conf.d/00-keyboard.conf`. The i3 `setxkbmap` line is
removed in favour of this.

### Multi-monitor: tray/primary stranded after undocking (autorandr)

i3 pins workspaces to outputs by name (`.config/i3/config`), but i3 doesn't
arrange the outputs themselves — X does, via `xrandr`. A bare
`xrandr ... --primary` in the config works on login, but if you unplug the
external monitor mid-session X doesn't promote the laptop back to primary, so
the i3 tray (`tray_output primary`) ends up on a dead output.

`autorandr` (installed by `autorandr.sh`) fixes this: it saves a profile per
monitor combination — outputs, resolutions, positions, and which is primary —
and reapplies the matching one. It ships udev + systemd user hooks that fire on
hotplug, and i3 also runs `autorandr --change` on login/restart. After install,
save the profiles once (commands printed by the script):

```bash
xrandr --output DP-1 --primary --auto --right-of eDP-1
autorandr --save docked     # external present
autorandr --save mobile     # laptop only (unplug first)
```

Profiles live in `~/.config/autorandr/` and are host-specific (keyed on monitor
EDID), so they're deliberately **not** tracked in this repo.

### nvim-treesitter errors on launch: `ENOENT ... 'tree-sitter'`

The neovim config runs nvim-treesitter's **`main`** branch, which (unlike the old
`master` branch with bundled parsers) compiles each parser at install time by
shelling out to the `tree-sitter` CLI, which in turn invokes a C compiler. A
fresh machine has neither, so launching nvim spams:

```
[nvim-treesitter/install/<lang>] error: Error during "tree-sitter build": ... ENOENT ... 'tree-sitter'
```

Two missing pieces (both handled by `nvim.sh`):

- **`tree-sitter` CLI** — brew's `tree-sitter` formula installs only the
  *library* (`libtree-sitter`); the binary lives in the separate
  `tree-sitter-cli` formula.
- **C compiler** — `tree-sitter build` calls `cc`; install `build-essential`.

After both are installed, open nvim and run `:TSUpdate` once to compile parsers.

### i3 login session missing from the GDM gear/cog menu

On a Wayland-only Ubuntu (26.04 `resolute` and later), GNOME ships **without a
standalone Xorg server** — only Xwayland is present. i3 is an **X11** window
manager whose session file (`/usr/share/xsessions/i3.desktop`) is
`Type=XSession` and needs a real `Xorg` binary to launch.

GDM filters out X11 sessions it can't actually start. With no Xorg installed it
drops both i3 entries, leaving only "Ubuntu (Wayland)" — and when only one
session remains, **GDM hides the session-picker gear entirely**. So the cog
appears "missing" when really there's just nothing else to pick.

Fix (handled by `i3.sh`):

```bash
sudo apt install xserver-xorg xinit
sudo reboot
```

After reboot the gear returns on the password screen (bottom-right) and i3 is
selectable.

### sur5r i3 repo 404 on new Ubuntu releases

The i3 docs (https://i3wm.org/docs/repositories.html) tell you to add the sur5r
apt repo, but it often has no build for a brand-new Ubuntu codename yet
(`apt update` → `404 Not Found` for `dists/<codename>/Release`). Just install
i3 from Ubuntu's official `universe` repo instead (what `i3.sh` does) and remove
the broken source:

```bash
sudo rm /etc/apt/sources.list.d/sur5r-i3.list
```
