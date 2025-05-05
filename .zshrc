alias vim=nvim
alias vi=nvim

# ENVIRONMENT
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$XDG_CONFIG_HOME/local/share"
export XDG_CACHE_HOME="$XDG_CONFIG_HOME/cache"

# DOTFILES
export DOTFILES="$HOME/.dotfiles"

# NOTES
export ZK_NOTEBOOK_DIR=~/notes

# EDITOR
export EDITOR=nvim
export VISUAL=nvim

# HISTORY
export HISTFILE="$HOME/.zsh_history" # HISTORY FILEPATH
export HISTSIZE=10000 # maximum events for internal history
export SAVEHIST=10000 # maximum events in history file

# ALIASES
source $DOTFILES/zsh/aliases.zsh
source $DOTFILES/zsh/functions.zsh

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


# VI MODE
# bindkey -v
# export KEYTIMEOUT=1

# PROMPT/THEME
# source $DOTFILES/zsh/themes/oxide.zsh-theme
eval "$(oh-my-posh init zsh --config ~/.dotfiles/ohmyposh/zen.toml)"

# SYNTAX HIGHLIGHTING
# source $DOTFILES/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

# ZSH HISTORY SUBSTRING SEARCH
# NB! Load this after syntax highlighting
# source $DOTFILES/zsh/plugins/zsh-history-substring-search/zsh-history-substring-search.zsh
# bindkey '^[[A' history-substring-search-up
# bindkey '^[[B' history-substring-search-down

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
export PATH=${PATH}:$HOME/.deno/bin

# Google cloud sdk components
# source "$(brew --prefix)/share/google-cloud-sdk/path.zsh.inc"
# export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"

# pnpm
export PNPM_HOME="/Users/sib/.config/local/share/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac
# pnpm end

