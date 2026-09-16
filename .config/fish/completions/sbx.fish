# Workspace names come from sbx's configured root, including stopped workspaces.
complete -c sbx -n '__fish_use_subcommand' -f -a new -d 'Create a workspace'
complete -c sbx -n '__fish_use_subcommand' -f -a open -d 'Open a workspace'
complete -c sbx -n '__fish_use_subcommand' -f -a delete -d 'Delete a workspace'
complete -c sbx -n '__fish_use_subcommand' -f -a list -d 'List workspaces'
complete -c sbx -n '__fish_use_subcommand' -f -a 'add project split tui tmux-reload doctor native help'

complete -c sbx -n '__fish_seen_subcommand_from open delete; and test (count (commandline -opc | string match -v -- "-*")) -eq 2' -f -a '(command sbx list --names 2>/dev/null)'
complete -c sbx -n '__fish_seen_subcommand_from open delete list' -f
complete -c sbx -n '__fish_seen_subcommand_from delete' -l yes -d 'Skip confirmation'
complete -c sbx -n '__fish_seen_subcommand_from list' -l names -d 'Print all managed workspace names'
