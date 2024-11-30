function pkill() {
  ps aux | fzf --height 40% --layout=reverse --prompt="Select process to kill: " | awk '{print $2}' | xargs -r sudo kill
}

# PROJECT: git-log
function logg() {
    git lg | fzf --ansi --no-sort \
        --preview 'echo {} | grep -o "[a-f0-9]\{7\}" | head -1 | xargs -I % git show % --color=always' \
        --preview-window=right:50%:wrap --height 100% \
        --bind 'enter:execute(echo {} | grep -o "[a-f0-9]\{7\}" | head -1 | xargs -I % sh -c "git show % | nvim -c \"setlocal buftype=nofile bufhidden=wipe noswapfile nowrap\" -c \"nnoremap <buffer> q :q!<CR>\" -")' \
        --bind 'ctrl-e:execute(echo {} | grep -o "[a-f0-9]\{7\}" | head -1 | xargs -I % sh -c "gh browse %")' \
}

function tskill() {
    tmux ls | fzf --height 40% --layout=reverse --prompt="Select tmux session to kill: " |  awk '{print substr($1, 1, length($1)-1)}' | xargs -r tmux kill-session -t
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
