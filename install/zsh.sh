if ! is-macos -o ! is-executable brew; then
  echo "Skipped: zsh"
  return
fi

brew install zsh
chsh -s /bin/zsh

echo "installing oh-my-zsh"
export ZSH="$HOME/.oh-my-zsh"; sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"

echo "installing fast-syntax-highlighting"
git clone https://github.com/zdharma/fast-syntax-highlighting.git "~/.oh-my-zsh/custom/plugins"

echo "installing fzf"
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
