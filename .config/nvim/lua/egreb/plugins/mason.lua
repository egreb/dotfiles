vim.pack.add({
	"https://github.com/mason-org/mason.nvim",
})

-- Installs language servers, linters and formatters into ~/.local/share/nvim/mason
-- and puts their binaries on $PATH (so `gopls`, `oxlint` etc. are found by the LSP config).
require("mason").setup()
