if is-macos; then
	echo "is mac, skipping snap"
	return
fi

if is-executable snap; then
	echo "snap is already installed.."
	return
fi

sudo apt update
sudo apt install snapd

if ! is-executable snap; then
	echo "failed to install snap"
	return
fi

apps=(
	slack
	spotify
	hugo
	chromium
	insomnia
	htop
)

for app in ${apps[@]}; do
	echo "installing $app..."
	snap install $app
done
