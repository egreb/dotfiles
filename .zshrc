
ZSH=$HOME/.oh-my-zsh
ZSH_THEME="avit"

alias dev="~/dev"
alias gri="git ls-files --ignored --exclude-standard | xargs -0 git rm -r" # fucked up .gitignore? run gri
alias fix="git diff --name-only | uniq | xargs code" # open code to merge conflicts
alias gs="git status"
alias reload="source ~/.zshrc"

# brew install tree
function t() {
  # Defaults to 3 levels deep, do more with `t 5` or `t 1`
  # pass additional args after
  tree -I '.git|node_modules|bower_components|.DS_Store' --dirsfirst --filelimit 15 -L ${1:-3} -aC $2
}

# Make utilities available
PATH="$DOTFILES_DIR/bin:$PATH"

# Path to your oh-my-zsh installation.
export ZSH=~/.oh-my-zsh

# Add wisely, as too many plugins slow down shell startup.
plugins=(
  git
  fast-syntax-highlighting
  node
  golang
  npm
  zsh-syntax-highlighting
  zsh-autosuggestions
)

source $ZSH/oh-my-zsh.sh

# FZF
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh
