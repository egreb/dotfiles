return {
	'stevearc/oil.nvim',
	---@module 'oil'
	---@type oil.SetupOpts
	dependencies = { { "echasnovski/mini.icons", opts = {} } },
	lazy = false,
	config = function()
		require('oil').setup({
			view_options = {
				show_hidden = true,
			},
			float = {
				padding = 2,
				max_width = 80,
				max_height = 40,
				border = 'rounded',
			},
		})
		-- file navigation
		vim.keymap.set("n", "<leader>b", ":Oil --float<CR>", { desc = "File navigation" })
		vim.keymap.set("n", "-", "<CMD>Oil<CR>", { desc = "Open parent directory" })
	end
}
