#!/usr/bin/env bash

# Get current dir
export DOTFILES_DIR DOTFILES_CACHE DOTFILES_EXTRA_DIR
DOTFILES_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOTFILES_CACHE="$DOTFILES_DIR/.cache.sh"
DOTFILES_EXTRA_DIR="$HOME/.extra"

# make utilities available
PATH="$DOTFILES_DIR/bin:$PATH"

# update dotfiles itself
echo "update dotfiles repo"

if [ -d "$DOTFILES_DIR/.git" ] 
then
	git --work-tree="$DOTFILES_DIR" --git-dir="$DOTFILES_DIR/.git" pull origin master
	echo "repo updated";
fi

if ! is-dir ~/.config
then 
	echo ".config does not exist, creating folder"
	mkdir ~/.config
fi

# Bunch of symlinks
echo "Install symlinks.."
ln -sfv "$DOTFILES_DIR/.zshrc" ~
ln -sfv "$DOTFILES_DIR/git/.gitconfig" ~
ln -sfv "$DOTFILES_DIR/git/.gitignore_global" ~
ln -sfv "$DOTFILES_DIR/.config/nvim" ~/.config
ln -sfv "$DOTFILES_DIR/.config/kitty" ~/.config
ln -sfv "$DOTFILES_DIR/.tmux.conf" ~
ln -sfv "$DOTFILES_DIR/.tmux.conf.osx" ~

if ! is-executable brew 
then 
	echo "installing brew..."
fi

# make sure brew is installed
if ! is-executable brew 
then
	echo "aborting, 'brew' command was not found.."
fi

