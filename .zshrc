# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi # Path to your oh-my-zsh installation.


# Path to your oh-my-zsh installation.
export ZSH=~/.oh-my-zsh
ZSH_THEME="powerlevel10k/powerlevel10k"

# Add wisely, as too many plugins slow down shell startup.
plugins=(
  git
  zsh-interactive-cd
)

source $ZSH/oh-my-zsh.sh

DISABLE_AUTO_TITLE=true

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

eval "$(fnm env --use-on-cd)"
source ~/.fzf.zsh

# ALIAS
alias t="tmux"
alias tk="tmux kill-session -t"
alias tl="tmux list-sessions"
alias ta="tmux attach -t"
alias tn="tmux new -s"

alias gs="git status"
alias gst="git status"
alias gp="git push"
alias reload="source ~/.zshrc"
alias vim="nvim"
# easy switch worktree and change tmux window name
function swt() {
  foo=$(git worktree list | fzf)
  dir=$(echo $foo | awk '{split($0,a); print $1}')
  branch=$(echo $foo | awk '{split($0,a); print $3}')

  tmux new-window -n $branch -c $dir
}

alias l="exa -l --icons --git -a"
alias lt="exa --tree --level=2 --long --icons --git"

# Make .dotfiles utilities available
PATH="$HOME/.dotfiles/bin:$PATH"

export EDITOR=nvim
export editor=nvim

# GO
export GOPATH=~/go
export PATH="$PATH:/usr/local/go/bin/"
export PATH="$PATH:$(go env GOPATH)/bin"

# npm
NPM_PACKAGES="${HOME}/.npm-packages"
export PATH="$PATH:$NPM_PACKAGES/bin"

