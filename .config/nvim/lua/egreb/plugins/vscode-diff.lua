return {
	'esmuellert/vscode-diff.nvim',
	dependencies = { 'MunifTanjim/nui.nvim' },
	cmd = 'CodeDiff',
	keys = {
		{
			'<leader>cd',
			'<cmd>:CodeDiff<cr>',
		},
	},
}
