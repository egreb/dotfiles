return {
	{
		'williamboman/mason-lspconfig.nvim',
		opts = {
			-- list of servers for mason to install
			ensure_installed = {
				'ts_ls',
				'html',
				'tailwindcss',
				'lua_ls',
				'emmet_ls',
				'jsonls',
				'bashls',
				'vue_ls',
			},
		},
		dependencies = {
			{
				'williamboman/mason.nvim',
				opts = {
					ui = {
						icons = {
							package_installed = '✓',
							package_pending = '➜',
							package_uninstalled = '✗',
						},
					},
				},
			},
			{
				'neovim/nvim-lspconfig',
			},
		},
	},
	{
		'WhoIsSethDaniel/mason-tool-installer.nvim',
		opts = {
			ensure_installed = {
				'prettier', -- prettier formatter
				'stylua', -- lua formatter
			},
		},
		dependencies = {
			'williamboman/mason.nvim',
		},
	},
}
