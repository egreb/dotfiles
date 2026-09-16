eval (/opt/homebrew/bin/brew shellenv)

set -U fish_greeting # disable fish greeting
set -U fish_key_bindings fish_vi_key_bindings
set -gx NVIM_APPNAME nvim-v2 # use the v2 config for Neovim and child tools
set -Ux EDITOR nvim # set Neovim as the default editor
set -Ux ZK_NOTEBOOK_DIR ~/notes
set -Ux TMUX_SESSIONIZER_PATHS "$HOME/Developer/work/agents $HOME/Developer/work/code $HOME/dotfiles"
set -Ux TMUX_SESSIONIZER_DEPTH 1
set -U DELTA_FEATURES "diff-so-fancy"
set -gx LG_CONFIG_FILE "$HOME/.config/lazygit/solarized.yml"
if test -f "$HOME/.config/lazygit/config.yml"
    set -gx LG_CONFIG_FILE "$HOME/.config/lazygit/config.yml,$LG_CONFIG_FILE"
end

fish_config theme choose solarized-classic

# alias
alias cc='sbx run --cpus 3 --memory 4g claude .'
alias claude='sbx run --cpus 3 --memory 4g claude .'
alias nvmi='NVIM_APPNAME=nvim-v2 nvim'

# abbreviations
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

bind -M insert \cf "~/dotfiles/.config/tmux/tmux-sessionizer.sh"
bind -M default \cf "~/dotfiles/.config/tmux/tmux-sessionizer.sh"

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

# opencode
fish_add_path /Users/sib/.opencode/bin

fish_add_path --path --move --prepend ~/.config/sbx/bin
