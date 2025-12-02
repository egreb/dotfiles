eval "$(starship init zsh)"
eval "$(zoxide init zsh)"

alias vim=nvim
alias vi=nvim
alias gcbrun="gh pr comment -b '/gcbrun'"

# NOTES
export ZK_NOTEBOOK_DIR=~/notes

# EDITOR
export EDITOR=nvim

# HISTORY
export HISTFILE=~/.history # HISTORY FILEPATH


# Set up fzf key bindings and fuzzy completion
source <(fzf --zsh)

bindkey -s ^f "tmux-sessionizer\n"

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:~/dotfiles/bin

eval "$(fnm env --use-on-cd)"
