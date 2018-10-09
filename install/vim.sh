if is-macos; then
    brew install vim
else
    if ! is-executable git; then
        echo "Skipped installing vim, missing git"
        exit 0
    fi

    mkdir -p "$DOTFILES_DIR/vim/tmp"
    git clone https://github.com/vim/vim.git "$DOTFILES_DIR/vim/tmp/"
    cd "$DOTFILES_DIR/vim/tmp/src" && make
    rm -rf "$DOTFILES_DIR/vim/tmp"
fi

# install plug(in) manager
curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
