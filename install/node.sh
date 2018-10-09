if is-macos -o; then
	if ! is-executable brew -o; then
		echo "brew missing, skipping node"
		return
	fi

	if ! is-executable node; then
		echo "installing node"
		brew install node
	else
		brew upgrade node
	fi
else
	if ! is-executable snap; then
		echo "snap missing, skipping node"
		return
	fi

	if ! is-executable node; then
		echo "installing node"
		snap install node
	else
		snap update node
	fi
fi

if ! is-executable npm -v; then
	echo "missing npm, skipping installing global packages"
	return
fi

apps=(
	tldr # manpage thing
	n # node version managment
)

if [ -d ~/.npm-global ]; then
	echo "creating global package folder"
	 mkdir ~/.npm-global
fi

echo "update npm config"
npm config set prefix "$HOME/.npm-global"

echo "updating npm."
npm install -g npm

echo "installing global packages..."
for app in ${apps[@]}; do
	echo "installing $app..."
	npm install -g $app
done

echo "all done"
