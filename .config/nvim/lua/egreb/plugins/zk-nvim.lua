-- note taking system
return {
	'zk-org/zk-nvim',
	config = function()
		require('zk').setup {
			lsp = {
				-- `config` is passed to `vim.lsp.start_client(config)`
				config = {
					cmd = { 'zk', 'lsp' },
					name = 'zk',
				},
				auto_attach = {
					enabled = true,
					filetypes = { 'markdown' },
				},
			},
		}
	end,
}
