if ! is-executable brew -o ! is-executable git; then
	echo "skipped node, missing brew and/or git"
	return
fi

apps=(
  node
  nvm
)

brew install "${apps[@]}"

export DOTFILES_BREW_PREFIX_NVM=`brew --prefix nvm`
set-config "DOTFILES_BREW_PREFIX_NVM" "$DOTFILES_BREW_PREFIX_NVM" "$DOTFILES_CACHE"

. "${DOTFILES_DIR}/system/.nvm"
nvm install 10
nvm alias default 10

# Global install with npm
packages=(
	npm
	tldr
)

npm install -g "${packages[@]}"
