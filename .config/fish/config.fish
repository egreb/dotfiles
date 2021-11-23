alias st="git status"
alias vim="nvim"

# Fish syntax highlighting
set -g fish_color_autosuggestion abb2bf 145 7
set -g fish_color_command  699959 71  2
set -g fish_color_match    699959 71  2

starship init fish | source

set NPM_PACKAGES "$HOME/.npm-packages"
set PATH $PATH $NPM_PACKAGES/bin
set MANPATH $NPM_PACKAGES/share/man $MANPATH

set PATH $PATH /home/sib/go/bin
