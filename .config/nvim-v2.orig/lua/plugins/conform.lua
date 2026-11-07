-- plugins
vim.pack.add({
	"https://github.com/stevearc/conform.nvim",
	"https://github.com/windwp/nvim-ts-autotag",
})

-- options
require("conform").setup({
	format_on_save = {
		timeout_ms = 8000,
		lsp_format = "fallback",
	},
	formatters_by_ft = {
		lua = { "stylua" },
		javascript = { "oxfmt" },
		javascriptreact = { "oxfmt" },
		typescript = { "oxfmt" },
		typescriptreact = { "oxfmt" },
		graphql = { "prettier" },
		go = { "goimports", "gofmt" },
		json = { "oxfmt" },
		sql = { "sql_formatter" },
	},
	formatters = {
		sql_formatter = {
			prepend_args = { "--language", "postgresql" },
		},
		prettier = {
			condition = function(_, ctx)
				return vim.fs.root(ctx.filename, { "biome.json" }) == nil
			end,
		},
	},
})

require("nvim-ts-autotag").setup()
