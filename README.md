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
Daily notes stay plain Markdown; render-markdown handles the display.
