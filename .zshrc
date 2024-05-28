# ENVIRONMENT
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$XDG_CONFIG_HOME/local/share"
export XDG_CACHE_HOME="$XDG_CONFIG_HOME/cache"

# DOTFILES
export DOTFILES="$HOME/.dotfiles"

# EDITOR
export EDITOR="nvim"
export VISUAL="nvim"

# HISTORY
export HISTFILE="$HOME/.zsh_history" # HISTORY FILEPATH
export HISTSIZE=10000 # maximum events for internal history
export SAVEHIST=10000 # maximum events in history file

# ALIASES
source $DOTFILES/zsh/aliases.zsh

# +------------+
# | NAVIGATION |
# +------------+
setopt AUTO_CD              # Go to folder path without using cd.

setopt AUTO_PUSHD           # Push the old directory onto the stack on cd.
setopt PUSHD_IGNORE_DUPS    # Do not store duplicates in the stack.
setopt PUSHD_SILENT         # Do not print the directory stack after pushd or popd.

setopt CORRECT              # Spelling correction
setopt CDABLE_VARS          # Change directory to a path stored in a variable.
setopt EXTENDED_GLOB        # Use extended globbing syntax.

# +---------+
# | HISTORY |
# +---------+
setopt EXTENDED_HISTORY          # Write the history file in the ':start:elapsed;command' format.
setopt SHARE_HISTORY             # Share history between all sessions.
setopt HIST_EXPIRE_DUPS_FIRST    # Expire a duplicate event first when trimming history.
setopt HIST_IGNORE_DUPS          # Do not record an event that was just recorded again.
setopt HIST_IGNORE_ALL_DUPS      # Delete an old recorded event if a new event is a duplicate.
setopt HIST_FIND_NO_DUPS         # Do not display a previously found event.
setopt HIST_IGNORE_SPACE         # Do not record an event starting with a space.
setopt HIST_SAVE_NO_DUPS         # Do not write a duplicate event to the history file.
setopt HIST_VERIFY               # Do not execute immediately upon history expansion.

# ENABLE SHELL AUTOCOMPLETION
autoload -U compinit; compinit
_comp_options+=(globdots) # With hidden files
# source $DOTFILES/zsh/completions.zsh
fpath=($DOTFILES/zsh/completions/src $fpath)

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

# VI MODE
bindkey -v
# export KEYTIMEOUT=1

# PROMPT/THEME
source $DOTFILES/zsh/themes/oxide.zsh-theme

# SYNTAX HIGHLIGHTING
source $DOTFILES/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

# ZSH HISTORY SUBSTRING SEARCH
# NB! Load this after syntax highlighting
source $DOTFILES/zsh/plugins/zsh-history-substring-search/zsh-history-substring-search.zsh
bindkey '^[[A' history-substring-search-up
bindkey '^[[B' history-substring-search-down

# INTERACTIVE CD
source $DOTFILES/zsh/plugins/zsh-interactive-cd.zsh

# FZF
source $DOTFILES/zsh/fzf.zsh

# NODE ENV MANAGER
eval "$(fnm env --use-on-cd)"

# PATH
export PATH=${PATH}:$HOME/go/bin
export PATH=${PATH}:$HOME/.bin
export PATH=${PATH}:$HOME/.dotfiles/bin

# Google cloud sdk components
source "$(brew --prefix)/share/google-cloud-sdk/path.zsh.inc"

