-- Enable the language servers. Their definitions (cmd, filetypes, root markers, settings)
-- live in `after/lsp/<name>.lua` and are picked up by the builtin `vim.lsp.config`.
vim.lsp.enable({
	"gopls", -- golang
	"tsc", -- typescript (TypeScript 7 native language server, the successor of tsgo)
	"lua_ls", -- lua
	"marksman", -- markdown
	"jsonls", -- json
	"oxlint", -- typescript / javascript linting
})
