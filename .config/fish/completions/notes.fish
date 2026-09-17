complete -c notes -f
complete -c notes -n 'not __fish_seen_subcommand_from today yesterday search new' -a today -d "Open today's note"
complete -c notes -n 'not __fish_seen_subcommand_from today yesterday search new' -a yesterday -d "Open yesterday's note"
complete -c notes -n 'not __fish_seen_subcommand_from today yesterday search new' -a search -d 'Search notes'
complete -c notes -n 'not __fish_seen_subcommand_from today yesterday search new' -a new -d 'Create a note'
complete -c notes -s h -l help -d 'Show usage'
