vim.pack.add({
	"https://github.com/folke/snacks.nvim",
})

-- LSP pickers open at the bottom of the screen (list on the left, preview on the right)
local bottom = { layout = { preset = "ivy" } }

require("snacks").setup({
	-- File/grep/buffer pickers etc. The default layout shows a preview of the selected item.
	picker = {
		enabled = true,
		sources = {
			-- Include dotfiles and dot directories (e.g. .cloudbuild). Gitignored files are still excluded.
			files = { hidden = true },
			lsp_definitions = bottom,
			lsp_implementations = bottom,
			lsp_references = bottom,
			lsp_symbols = bottom,
			diagnostics = bottom,
		},
	},
})
