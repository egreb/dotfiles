return {
	'saghen/blink.cmp',
	version = '1.*',
	event = { 'BufReadPre', 'BufNewFile' },
	dependencies = {
		{ 'antosha417/nvim-lsp-file-operations', config = true },
		{ 'folke/lazydev.nvim',                  opts = {} },
	},
	---@module 'blink.cmp'
	---@type blink.cmp.Config
	opts = {
		appearance = {
			-- 'mono' (default) for 'Nerd Font Mono' or 'normal' for 'Nerd Font'
			-- Adjusts spacing to ensure icons are aligned
			nerd_font_variant = 'mono',
		},
		fuzzy = { implementation = 'prefer_rust_with_warning' },
		signature = { window = { border = 'padded' } },
		completion = {
			accept = { auto_brackets = { enabled = false } },
			documentation = {
				auto_show = true,
				auto_show_delay_ms = 250,
				window = {
					border = 'padded',
				},
			},
			menu = {
				border = 'padded',
				draw = {
					gap = 2,
					columns = { { 'label', 'label_description', gap = 1 }, { 'kind_icon', 'kind' } },
				},
			},
			-- Display a preview of the selected item on the current line
			ghost_text = { enabled = false },
		},
		sources = {
			default = { 'lsp', 'path', 'buffer' },
		},
	},
}
