-- Markdown parsers and mini.icons are provided by the existing plugins.
vim.pack.add({ "https://github.com/MeanderingProgrammer/render-markdown.nvim" })

require("render-markdown").setup({
	completions = { lsp = { enabled = true } },
	-- Keep the layout rendered when entering character, line, or block Visual mode.
	render_modes = { "n", "c", "t", "v", "V", "\22" },
	-- Revealing markup under the cursor / selection shifts rendered content.
	anti_conceal = { enabled = false },
	win_options = {
		concealcursor = { rendered = "nvc" },
	},
})

vim.keymap.set("n", "<leader>um", "<cmd>RenderMarkdown toggle<CR>", { desc = "Toggle Markdown rendering" })
