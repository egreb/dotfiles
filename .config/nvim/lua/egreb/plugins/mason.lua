vim.pack.add({
	"https://github.com/mason-org/mason.nvim",
})

-- Installs language servers, linters and formatters into ~/.local/share/nvim/mason
-- and puts their binaries on $PATH (so `gopls`, `oxlint` etc. are found by the LSP config).
require("mason").setup()

-- Mason package names (not lspconfig names) of everything the config relies on.
local tools = {
	-- LSPs
	"gopls", -- golang
	"tsc", -- typescript (TypeScript 7 native language server)
	"lua-language-server", -- lua
	"marksman", -- markdown
	"json-lsp", -- json
	"oxlint", -- typescript / javascript linting
	-- Formatters (gofmt ships with Go itself)
	"goimports", -- golang
	"oxfmt", -- typescript / json
	"stylua", -- lua
}

-- Make sure everything above is installed on startup. Runs asynchronously so it doesn't block.
local registry = require("mason-registry")
registry.refresh(function()
	for _, name in ipairs(tools) do
		local ok, pkg = pcall(registry.get_package, name)
		if not ok then
			vim.schedule(function()
				vim.notify(("mason: unknown package %q"):format(name), vim.log.levels.WARN)
			end)
		elseif not pkg:is_installed() then
			pkg:install()
		end
	end
end)
