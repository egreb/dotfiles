# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

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
alias gcbrun="gh pr comment -b '/gcbrun'"
alias notify="noti gh pr checks --watch"
alias 1c="git merge origin/main && git reset --soft origin/main"
alias ts="tmux new-session -t $1"

function squash() {
  main_branch=$(basename $(git symbolic-ref --short refs/remotes/origin/HEAD))
  git merge origin/$main_branch && git reset --soft origin/$main_branch
}

# switch tmux sessions using fzf
function ss() {
  if [ -z "$TMUX" ]; then
    tmux attach-session -t "$(tmux ls -F '#S' | fzf --prompt='Select session: ')"
  else
    tmux switch-client -t "$(tmux ls -F '#S' | fzf --prompt='Select session: ')"
  fi
}

# easy switch worktree and change tmux window name
function swt() {
  selected=$(git worktree list | fzf)
  dir=$(echo $selected | awk '{split($0,a); print $1}')
  branch=$(echo $selected | awk '{split($0,a); print $3}')

  if ! [ -z "$dir" ]; then
    tmux new-window -n $branch -c $dir
  fi
}

# easy switch worktree and change tmux window name
# zle -N swt
# bindkey '^W+l' swt

function gwa() {
  BRANCH=$1

  git worktree add $1 -b $1 >/dev/null 2>&1
  worktree=$(git worktree list | fzf --query="$BRANCH" -1)
  dir=$(echo $worktree | awk '{split($0,a); print $1}')

  if tmux info &> /dev/null; then
    tmux new-window -n $BRANCH -c $dir
  else
    cd $dir
  fi
}

alias l="exa -l --icons --git -a"
alias lt="exa --tree --level=2 --long --icons --git"
alias pm="pnpm"
alias vpr="gh pr view -w"
alias cpr="gh pr create -w"

# Make .dotfiles utilities available
PATH="$HOME/.dotfiles/bin:$PATH"

export EDITOR=nvim
export editor=nvim

# GO
export GOPATH=~/go
export PATH="$PATH:$(go env GOPATH)/bin"
export PATH="$PATH:$(go env GOBIN)"

# npm
NPM_PACKAGES="${HOME}/.npm-packages"
export PATH="$PATH:$NPM_PACKAGES/bin"

PATH="$HOME/.bin:$PATH"

# gcloud
source "$(brew --prefix)/share/google-cloud-sdk/path.zsh.inc"

# syntax highlighting fish style
source /opt/homebrew/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
