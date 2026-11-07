return {
	'folke/snacks.nvim',
	dependencies = {
		'echasnovski/mini.icons',
	},
	priority = 1000,
	lazy = false,
	opts = {
		dim = { enabled = true },
		indent = { enabled = true },
		scroll = { enabled = true },
		notifier = { enabled = true },
		picker = {
			enabled = true,
			hidden = false,
			ignored = false,
			preview = false,
			cwd = vim.fn.getcwd(),
			preset = 'ivy',
			sources = {
				explorer = {
					layout = {
						layout = {
							position = 'right',
						},
					},
				},
			},
		},
		explorer = {
			enabled = true,
		},
		input = {
			enabled = true,
			icon = ' ',
			icon_hl = 'SnacksInputIcon',
			icon_pos = 'left',
			prompt_pos = 'title',
			win = { style = 'input' },
			expand = true,
		},
		scratch = { enabled = true },
		lazygit = {
			enabled = true,
			configure = true,
			win = {
				style = 'lazygit',
			},
			selectedLineBgColor = { bg = '#333' },
		},
		gitbrowse = { enabled = true },
		bufdelete = {},
	},

	keys = {
		{
			'<C-n>',
			function()
				Snacks.explorer()
			end,
			desc = 'Explorer',
		},
		{
			'<leader>nn',
			function()
				Snacks.scratch()
			end,
			desc = 'Toggle Scratch Buffer',
		},
		{
			'<leader>ns',
			function()
				Snacks.scratch.select()
			end,
			desc = 'Select Scratch Buffer',
		},
		{
			'<leader>lg',
			function(opts)
				Snacks.lazygit.open(opts)
			end,
			desc = '[L]azy[G]it',
		},
		{
			'<leader>gb',
			function(opts)
				Snacks.git.blame_line(opts)
			end,
			'[G]it [B]lame',
		},
		{
			'<leader>go',
			function(opts)
				Snacks.gitbrowse.open(opts)
			end,
			'[G]it [B]lame',
		},
		{
			'<leader>tt',
			function(opts)
				Snacks.terminal.toggle(_, opts)
			end,
			'[T]oggle [T]erminal',
		},
	},
}
