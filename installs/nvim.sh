set -e

echo "installing neovim"
brew install neovim

# nvim-treesitter is on its `main` branch, which compiles parsers at startup by
# shelling out to the `tree-sitter` CLI plus a C compiler. Without these you get
# `Error during "tree-sitter build": ... ENOENT ... 'tree-sitter'` on launch.
# brew's `tree-sitter` formula ships only the library, so the CLI is a separate
# `tree-sitter-cli` formula.
echo "installing tree-sitter CLI (needed by nvim-treesitter main branch)"
brew install tree-sitter-cli

echo "installing build-essential (provides cc/gcc for tree-sitter to compile parsers)"
sudo apt install -y build-essential

echo "done. open nvim and run :TSUpdate to compile parsers."
