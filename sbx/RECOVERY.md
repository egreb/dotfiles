# Recovery record — 2026-09-17

This is a tested reconstruction, not a byte-for-byte restoration. The deleted
nested Git repository and its commit history were not recovered.

Source evidence: saved session logs from 2026-09-15, preserved separately at
`/Users/sib/sbx-recovery-20260917/evidence`. Historical source fragments and the
extraction/replay scripts also remain in that recovery directory.

## Recovered source

- App dispatcher and native-executable discovery.
- Workspace creation, agent/prompt selection, clone progress, and later form edits.
- Inventory, preview sanitization, and activation logic.
- Workspace path checks, Git clone/bootstrap logic, and Node-version handling.
- Core tmux wrappers, project-session metadata checks, and tool-window setup.
- Workspace form, clone-progress renderer, Claude theme initialization.
- Configuration loader and original tests present in the logs.

Recorded text transformations were replayed where their inputs were available.
Imports and integration points were repaired during compilation.

## Reconstructed code

Files named `*_reconstructed.go`, `internal/config/reconstructed.go`,
`internal/model/model.go`, and `internal/tui/tui.go` supply missing lifecycle,
configuration parsing, tmux overlay/sidecar behavior, data types, and picker code.
Process helper functions, the entry point, build files, examples, and documentation
were recreated. Reconstructed tests supplement the original recovered tests.
The original full integration suite and project manifest were not recovered.
Use the existing host configuration and manifest, not the generic examples.

## Verification

`make check` runs Go tests, builds the executable, runs go vet, and exercises an
isolated tmux server with local Git repositories and a fake Docker CLI. It covers
creation, independent clones, session metadata, tool windows, sidecars, key
forwarding, reuse, passive listing, unsafe deletion rejection, rollback, and exact
workspace deletion. The integration suite also drives prefix+f through an attached
tmux client, verifying popup context and selection from agent and project panes. Layout tests cover long names and narrow terminals.

This does not verify a real Docker Sandbox, host authentication, actual developer
tools, or the installed macOS tmux version. The recovery environment is Linux.
A macOS cross-build verifies compilation only; build/install on the host for use.
