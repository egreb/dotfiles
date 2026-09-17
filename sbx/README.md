# sbx workspace orchestrator

Host-side Go wrapper for Docker Sandboxes, with multi-project workspaces and a
tmux session picker. This source was reconstructed after accidental deletion;
see [RECOVERY.md](RECOVERY.md) for provenance and verification limits.

## Build and install

Requires Go 1.25+, tmux with popup support, Git, and Docker's native `sbx` CLI.
Fish is the default host shell. Node projects use fnm and pnpm; Go projects use Go.
The picker is built in; fzf is used for non-form project selection.

```sh
make check
make install
sbx tmux-reload
```

`make install` builds for the machine running it and atomically replaces
`~/.config/sbx/bin/sbx`. It does not overwrite your configuration, project manifest,
workspaces, or native Docker executable. Keep that directory first in PATH.
Existing `~/.config/sbx/config.yaml` and `projects.yaml` should continue to be used.
The example YAML files are for new installations, not replacements for your manifest.

## Usage

```sh
sbx new                         # form: name, agent, prompt, repository selection
sbx new --all --no-open example
sbx new --agent codex --prompt 'Fix the failing tests' example
sbx open                        # all-workspace picker
sbx open example
sbx project backend
sbx add frontend
sbx split frontend
sbx list --names
sbx delete example              # requires typing the exact workspace name
sbx doctor
```

Unrecognized commands forward to Docker's native CLI. `sbx native COMMAND` does
so explicitly. The native executable must differ from this wrapper.

The tool creates one Docker Sandbox per workspace. Project shells, Neovim,
lazygit, and hunk run on the host. Project sessions have three windows:
Neovim v2, lazygit, and `hunk diff`. Exiting a tool leaves its shell available.
Repository mirrors live under `~/.cache/sbx/git`; workspace clones are independent.
Node versions are resolved by fnm before configured bootstrap hooks and pnpm install.

## tmux

The dedicated server defaults to `sbx-v1`, using your configured tmux prefix.

| Key after prefix | Action |
| --- | --- |
| `i` | New workspace form |
| `o` | All workspace/project sessions |
| `f` | Current workspace sessions |
| `s` | Open a session beside the current pane |
| `a` | Add a project to the current workspace |
| `n` | Cycle panes |
| `r` | Reload config and reapply sbx bindings |

Your separate `Shift-n` work-notes popup remains defined by dotfiles.

In the picker, type to filter; Ctrl-j/k selects; Enter opens; Tab or Ctrl-h/l
controls the preview; Ctrl-r refreshes; F1 shows help; Escape clears or closes.
The list receives 60% of wide popups, and the full selected workspace/project
appears above it. Narrow popups automatically hide the preview.
Sidecars attach to existing sessions; closing a sidecar leaves its target alive.

Only marked direct children of the configured workspace root can be deleted.
Passive inventory never calls Docker. Creation rollback removes only resources
reserved by that creation attempt.
