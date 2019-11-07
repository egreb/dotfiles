export LC_ALL=en_US.UTF-8
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

ZSH_THEME="frisk"

# Make utilities available
PATH="$DOTFILES_DIR/bin:$PATH"

# Path to your oh-my-zsh installation.
export ZSH=~/.oh-my-zsh

# history
HIST_STAMPS="hh:ss"
HISTFILE="$HOME/.zsh_history"
HISTSIZE=10000000
SAVEHIST=10000000
setopt HIST_IGNORE_ALL_DUPS

# Add wisely, as too many plugins slow down shell startup.
plugins=(
  git 
  fast-syntax-highlighting
)

source $ZSH/oh-my-zsh.sh

# Finally we can source the dotfiles (order matters)
for DOTFILE in "$DOTFILES_DIR"/system/.{alias,functions}; do
  [ -f "$DOTFILE" ] && . "$DOTFILE"
done

 #BINS
export PATH=$PATH:~/.bin

# FZF
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

# az autocompletion
autoload -U +X bashcompinit && bashcompinit
source /usr/local/Cellar/azure-cli/2.0.67/etc/bash_completion.d/az
