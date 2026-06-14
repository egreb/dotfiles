# Dotfiles

Personal dotfiles. Configs live under `.config/<app>/` in this repo and are
symlinked into `~/.config/<app>` (e.g. `~/.config/fish -> ../dotfiles/.config/fish`).

## i3 keybindings cheat sheet

There is an HTML cheat sheet at `.config/i3/keybindings.html` that documents the
active i3 keybindings.

**Whenever you add, remove, or change a keybinding in `.config/i3/config`, update
`.config/i3/keybindings.html` in the same change so the two stay in sync.** The
config file is the source of truth; the HTML mirrors it.

Open the cheat sheet with `xdg-open ~/.config/i3/keybindings.html`.
