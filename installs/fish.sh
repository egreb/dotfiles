set -e

echo "installing fish"
brew install fish

echo "adding fish to /etc/shells"
sudo command -v fish | sudo tee -a /etc/shells

echo "setting fish as default shell"
sudo chsh -s "$(command -v fish)"

echo "finishing installing fish"
