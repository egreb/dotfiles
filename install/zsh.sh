if ! is-macos -o ! is-executable brew -o; then
  echo "Skipped: zsh"
  return
fi

brew install zsh
chsh -s /bin/zsh

if ! is-executable zsh; then
  echo "installing oh-my-zsh"
  export ZSH="$HOME/.oh-my-zsh"; sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
fi

echo "installing fast-syntax-highlighting"
git clone https://github.com/zdharma/fast-syntax-highlighting.git "~/.oh-my-zsh/custom/plugins"


if ! is-executable fzf --help; then
  echo "installing fzf"
  git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
  ~/.fzf/install
else
  echo "updating fzf"
  cd ~/.fzf && git pull && ./install
fi
