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
sbx delete                      # workspace picker, then exact-name confirmation
sbx delete example              # requires typing the exact workspace name
sbx doctor
```

Unrecognized commands forward to Docker's native CLI. `sbx native COMMAND` does
so explicitly. The native executable must differ from this wrapper.

The tool creates one Docker Sandbox per workspace. Project shells and Neovim
run on the host. Project sessions have two windows: Neovim v2 in tab 1 and an
empty shell in tab 2. Exiting Neovim leaves its shell available.
Repository mirrors live under `~/.cache/sbx/git`; workspace clones use Git transport and do not hardlink objects to the cache. `sbx new` and `sbx add` verify mirrors before using them and rebuild damaged mirrors from the configured upstream. Replaced mirrors are retained under `.PROJECT-damaged-*/mirror.git` in the cache directory for diagnosis; they can be removed once no longer needed. Cache operations are locked per project; if another command is using the same cache, retry when it finishes.

This repairs future clones only. Existing workspaces with missing objects need a fresh clone; preserve any local changes, commits, and stashes before replacing them.
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
| `Shift-D` | Choose a workspace to delete, then type its name to confirm |
| `n` | Cycle panes |
| `r` | Reload config and reapply sbx bindings |

Your separate `Shift-n` work-notes popup remains defined by dotfiles.

The picker groups sessions beneath workspace names, expanding the selected workspace.
Use `h/l` (or Left/Right) for previous/next workspace and `j/k` (or Up/Down,
Ctrl-j/k) for sessions within that workspace. Returning to a workspace restores
its selected session. Space folds the current group; moving within it expands it.
Press `/` to edit the filter, then Enter or Escape to return to navigation;
Enter opens the selected session. `s` toggles workspace-name and last-used sorting
(default: last used). Recency uses the latest tmux session attachment in each
workspace; unavailable history sorts last and is lost when those sessions are removed.
Tab or Ctrl-h/l controls the initially hidden preview; Ctrl-r refreshes; F1 shows
help; Escape clears a filter or closes the picker. Narrow popups hide the preview.
The full selected workspace/project appears above the list.
Press `D` (or Ctrl-d) on a workspace or any of its sessions to delete that whole
workspace. The picker names the workspace and warns that files, uncommitted
changes, the sandbox, and sessions will be removed; only `y` confirms.
After deletion, `y` opens the new-workspace form; `n`, Enter, or Escape opens the
most recently used remaining workspace. If none remain, it returns to a host
shell. Canceling the deletion itself keeps the selection in the picker.
The separate `sbx delete` command still requires exact-name confirmation.
Sidecars attach to existing sessions; closing a sidecar leaves its target alive.

Only marked direct children of the configured workspace root can be deleted.
Deletion sends SIGTERM to host processes in the workspace’s managed tmux panes,
waits up to two seconds, then sends SIGKILL to survivors before removing the
sandbox and workspace. Processes already detached from those process trees
are not discoverable by this cleanup.
Passive inventory never calls Docker. Creation rollback removes only resources
reserved by that creation attempt.

The delete popup also works from the sandbox agent pane: press your tmux prefix,
then Shift-D. It runs on the host and can delete the workspace you are currently
using; its sessions close after confirmation. Escape cancels the picker.
