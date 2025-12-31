return {
	'mcncl/alabaster.nvim',
	lazy = false,
	priority = 1000,
	config = function()
		require('alabaster').setup {
			style = 'dark', -- "light" or "dark"
		}
		vim.cmd.colorscheme 'alabaster'
	end,
}
