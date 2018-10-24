if ! is-macos -o ! is-executable brew -o; then
  echo "Missing brew, skipped: zsh"
  return
fi

if is-macos -o; then
  brew install zsh
  # use zsh installed by brew
  sudo -s 'echo /usr/local/bin/zsh >> /etc/shells' && chsh -s /usr/local/bin/zsh
else
  sudo apt install zsh
  chsh -s /usr/bin/zsh
fi

if ! is-executable zsh; then
  echo "skipping oh-my-zsh and plugins, missing zsh."
  return
fi

echo "installing oh-my-zsh"
export ZSH="$HOME/.oh-my-zsh"; sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"

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
