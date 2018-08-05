export PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"

# If not running interactively, don't do anything
[ -z "$PS1" ] && return

#READLINK=$(which greadlink || which readlink)
#CURRENT_SCRIPT=$BASH_SOURCE

# Resolve DOTFILES_DIR (assuming ~/.dotfiles on distros without readlink and/or $BASH_SOURCE/$0)
export PATH=$HOME/bin:/usr/local/bin:$PATH
if [[ -n $CURRENT_SCRIPT && -x "$READLINK" ]]; then
  SCRIPT_PATH=$($READLINK -f "$CURRENT_SCRIPT")
  DOTFILES_DIR=$(dirname "$(dirname "$SCRIPT_PATH")")
elif [ -d "$HOME/.dotfiles" ]; then
  DOTFILES_DIR="$HOME/.dotfiles"
else
  echo "Unable to find dotfiles, exiting."
  return
fi

ZSH_THEME="egreb"

# Make utilities available
PATH="$DOTFILES_DIR/bin:$PATH"

# Path to your oh-my-zsh installation.
export ZSH=~/.oh-my-zsh

HIST_STAMPS="hh:ss"

# Add wisely, as too many plugins slow down shell startup.
plugins=(git fast-syntax-highlighting)

source $ZSH/oh-my-zsh.sh
# Preferred editor for local and remote sessions
# if [[ -n $SSH_CONNECTION ]]; then
#   export EDITOR='vim'
# else
#   export EDITOR='mvim'
# fi

export JAVA_HOME=$(/usr/libexec/java_home)

#GO PATHS
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOPATH/bin

# BINS
export PATH=$PATH:~/.bins

# FZF
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

