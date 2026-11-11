-- Keep installed parsers in sync with the plugin version. Must be registered before `vim.pack.add`.
vim.api.nvim_create_autocmd("PackChanged", {
	desc = "Run :TSUpdate when nvim-treesitter is updated",
	callback = function(ev)
		local data = ev.data
		if data.spec.name == "nvim-treesitter" and data.kind == "update" then
			if not data.active then
				vim.cmd.packadd("nvim-treesitter")
			end
			vim.cmd("TSUpdate")
		end
	end,
})

vim.pack.add({
	"https://github.com/nvim-treesitter/nvim-treesitter",
})

-- Parsers to install (no-op for the ones already installed). Runs asynchronously.
-- Needs `tree-sitter-cli` (installed by Mason) and a C compiler.
require("nvim-treesitter").install({
	"go",
	"gomod",
	"gowork",
	"sql", -- inline sql in go files, see after/queries/go/injections.scm
	"typescript",
	"tsx",
	"javascript",
	"json",
	"lua",
	"markdown",
	"markdown_inline", -- needed by render-markdown
})

-- Treesitter highlighting is provided by Neovim, it just has to be started per buffer
vim.api.nvim_create_autocmd("FileType", {
	desc = "Start treesitter highlighting",
	callback = function()
		pcall(vim.treesitter.start)
	end,
})
