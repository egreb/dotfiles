if ! is-executable code ; then
	echo "code is not installed, install code then try again"
	return
fi

if [ -z $DOTFILES_DIR ]; then
	echo "DOTFILES_DIR variable not set, aborting"
	return
fi

if [ -d "$DOTFILES_DIR/vscode" ]; then
	echo "missing vscode src folder, aborting"
fi

SRC="$DOTFILES_DIR/vscode/User"

if is-macos; then
	DESTINATION="$HOME/Library/Application Support/Code/User"
else
	DESTINATION="$HOME/.config/Code/User"
fi

echo "copying vscode settings to $DESTINATION"

if [ -d "$DESTINATION" ]; then
	echo "deleting $DESTINATION"
	rm -rf $DESTINATION && "Deleted"
fi

echo "Linking $DESTINATION to $SRC"
ln -s "$SRC" "$DESTINATION" && echo "Linked!"

extensions=(
	ms-vscode.go # golang
	ms-vscode.sublime-keybindings # sublime keybindings
	peterjausovec.vscode-docker # docker files syntax
	zyst.egoist-one # theme
	dbaeumer.vscode-eslint # eslint
	pkief.material-icon-theme # icons
	zhuangtongfa.material-theme # theme
	esbenp.prettier-vscode # prettier
	ryannaddy.laravel-artisan # laravel artisan
	felixfbecker.php-debug # php debug
	neilbrayfield.php-docblocker # php docblocker
	bmewburn.vscode-intelephense-client # php intellisense
	kokororin.vscode-phpfmt # php formatting	
)

echo "installing extensions..."

extension_install() {
	local extension=$1
	echo "$extension"
	code --install-extension $extension
}

for EXTENSION in ${extensions[@]}; do
	extension_install $EXTENSION
done


echo "done copying code config!"
