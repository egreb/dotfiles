alias l="exa -l --icons --git -a"
alias lt="exa --tree --level=2 --long --icons --git"
alias pm="pnpm"
alias vpr="gh pr view -w"
alias cpr="gh pr create -w"

# GIT ALIASES
alias gs='git status'
alias ga='git add'
alias gp='git push'
alias gpo='git push origin'
alias gtd='git tag --delete'
alias gtdr='git tag --delete origin'
alias gr='git branch -r'
alias gplo='git pull origin'
alias gb='git branch '
alias gc='git commit'
alias gd='git diff'
alias gco='git checkout '
alias gl='git log'
alias gr='git remote'
alias grs='git remote show'
alias glo='git log --pretty="oneline"'
alias glol='git log --graph --oneline --decorate'
alias g="git"
alias open_repo="gh repo view --web \"\$(gh repo ls cybr-ai --json name,owner --limit 100 | jq -r '.[] | .owner.login + \"/\" + .name' | fzf --print0 --exit-0)\""
alias gcbrun="gh pr comment -b '/gcbrun'"
alias ta="tmux attach-session -t"
alias ts="tmux new-session -t"
alias ot="~/.dotfiles/bin/open-worktree.sh"

# notify when deploy is complete
alias notify="noti gh pr checks --watch"
