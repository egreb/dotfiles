vim.pack.add({
	"https://github.com/folke/snacks.nvim",
})

-- LSP pickers open at the bottom of the screen (list on the left, preview on the right)
local bottom = { layout = { preset = "ivy" } }

require("snacks").setup({
	-- Lazygit in a floating terminal, toggled with `<leader>lg`, see keymaps. Snacks gives it a theme
	-- matching the colorscheme, and editing a file from lazygit (`e`) opens it in this Neovim.
	lazygit = { enabled = true },

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
