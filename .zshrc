eval "$(starship init zsh)"
eval "$(zoxide init zsh)"

alias vim=nvim
alias vi=nvim
alias gcbrun="gh pr comment -b '/gcbrun'"
alias ll="exa -l"

setopt inc_append_history

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
. "/home/sib/.deno/env"
# Initialize zsh completions (added by deno install script)
autoload -Uz compinit
compinit
# bun completions
[ -s "/home/sib/.bun/_bun" ] && source "/home/sib/.bun/_bun"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# fnm
FNM_PATH="/home/sib/.local/share/fnm"
if [ -d "$FNM_PATH" ]; then
  	export PATH="$FNM_PATH:$PATH"
  # eval "`fnm env`"
	eval "$(fnm env --use-on-cd)"
fi
