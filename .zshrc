eval "$(starship init zsh)"
eval "$(zoxide init zsh)"

alias vim=nvim
alias vi=nvim

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

new_tmux () {
  session_dir=$(zoxide query --list | fzf)
  session_name=$(basename "$session_dir")

  if tmux has-session -t $session_name 2>/dev/null; then
    if [ -n "$TMUX" ]; then
      tmux switch-client -t "$session_name"
    else
      tmux attach -t "$session_name"
    fi
    notification="tmux attached to $session_name"
  else
    if [ -n "$TMUX" ]; then
      tmux new-session -d -c "$session_dir" -s "$session_name" && tmux switch-client -t "$session_name"
      notification="new tmux session INSIDE TMUX: $session_name"
    else
      tmux new-session -c "$session_dir" -s "$session_name"
      notification="new tmux session: $session_name"
    fi
  fi

  if [-s "$session_name" ]; then
    notify-send "$notification"
  fi
}

alias tm=new_tmux
