# Find the tmux pane a process belongs to by walking up its parent chain.
# Usage: whichpane <pid> [--go]   (--go switches to that pane)

function whichpane
    set -l pid $argv[1]
    if test -z "$pid"
        echo "usage: whichpane <pid> [--go]" >&2
        return 1
    end

    if not ps -p $pid >/dev/null
        echo "no process with pid $pid (it has exited)" >&2
        return 1
    end

    set -l panes (tmux list-panes -a -F '#{pane_pid} #{session_name}:#{window_index}.#{pane_index} #{pane_current_path}')
    while test -n "$pid"; and test "$pid" -gt 1
        ps -o pid=,ppid=,tty=,etime=,command= -p $pid | string trim | string shorten -m 120 | sed 's/^/  /' >&2
        for line in $panes
            set -l f (string split ' ' -- $line)
            if test "$f[1]" = "$pid"
                echo "→ $f[2]  $f[3]"
                if contains -- --go $argv
                    tmux switch-client -t $f[2]
                end
                return 0
            end
        end
        set pid (ps -o ppid= -p $pid | string trim)
    end

    echo "$argv[1] is not inside a tmux pane (chain above ends at launchd; ppid 1 right away = orphaned)" >&2
    return 1
end
