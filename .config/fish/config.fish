eval (/opt/homebrew/bin/brew shellenv)

set -U fish_greeting # disable fish greeting
set -U fish_key_bindings fish_vi_key_bindings
set -Ux EDITOR nvim # set nvim as default shell
set -Ux ZK_NOTEBOOK_DIR ~/notes
set -Ux CLAUDE_CONFIG_DIR ~/.claude-work
set -U DELTA_FEATURES "diff-so-fancy"

fish_config theme choose catppuccin-frappe

# alias
abbr ta "tmux attach-session -t"
abbr ts "tmux new-session -t"
abbr n "nvim"
abbr v "nvim"
abbr pm "pnpm"
abbr gcbrun "gh pr comment -b '/gcbrun'"
abbr notify "noti gh pr checks --watch"
abbr vpr "gh pr view -w"
abbr cpr "gh pr create -w"
abbr lg "lazygit"

# Set up fzf key bindings
fzf --fish | source

 # NODE ENV MANAGER
eval "$(fnm env --use-on-cd)"

bind -M insert \cf "tmux-sessionizer"
bind -M default \cf "tmux-sessionizer"

# uv
fish_add_path ~/.config/bin
fish_add_path ~/dotfiles/bin
fish_add_path ~/go/bin
fish_add_path "/Users/sib/.local/bin"
fish_add_path "/Users/sib/.config/herd-lite/bin"

function ws
    set selected (git worktree list | grep -v '(bare)' | awk '{print $1}' | fzf --prompt="Select worktree: ")
    
    if test -n "$selected"
        echo "Changing to: $selected"
        cd "$selected"
    else
        echo "No worktree selected"
    end
end


# Added by OrbStack: command-line tools and integration
# This won't be added again if you remove it.
source ~/.orbstack/shell/init2.fish 2>/dev/null || :
