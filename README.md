# my dotfiles

## Requirements

- Stow

## Setup

```bash
cd dotfiles_repo
stow .
```

## Notes from the terminal

In tmux, press **Ctrl-s, then Shift-n** to open the work-vault picker in a centered
Neovim popup. Use `:wq` to save and close it (or `:q` if unchanged).
Reload tmux configuration with **Ctrl-s, then r** after installing the binding.

In Rex, press **Cmd-Shift-O** (or run "Open Work Notes" from the command palette)
for the same picker in a floating pane, and **Cmd-Option-N** ("New Work Note") to
create a note. Both are defined in `.config/rex/init.lua`; run `rex config reload`
after changing it.

The `notes` command also opens the work vault (`~/vaults/work`) directly from
any directory. Fish already adds this repository's `bin` directory to PATH.

```bash
notes        # browse work notes
notes today  # open/create today's daily note
notes search # search work notes
notes new    # prompt for a new work note
```

Inside Neovim, use `<leader>mq` to browse notes, `<leader>ms` to search,
and `<leader>mt` for today's note.
Daily notes stay plain Markdown; md-render.nvim handles the display (notes open
as editable source, `<leader>mo` toggles the rendered read mode).

## Colors

Catppuccin Mocha is used throughout Ghostty, fish, fzf, both Neovim configs,
lazygit, Hunk, tuicr, tmux, Alacritty, bat, and delta. Neovim uses the
`catppuccin/nvim` theme plugin, which also styles the other plugins.
Ghostty, Hunk, tuicr, and bat use bundled themes; fish and lazygit use the
official Catppuccin presets. fzf shares `~/.config/fzf/config` between fish
and the tmux pickers; use a recent fzf with `FZF_DEFAULT_OPTS_FILE` support.

After stowing new config files, reload Ghostty, start a fresh fish shell,
restart Neovim and other tools, and reload tmux with **Ctrl-s, then r**.
Use recent bat/delta releases that include the `Catppuccin Mocha` syntax theme.
