vim.pack.add({
	"https://github.com/stevearc/conform.nvim",
})

require("conform").setup({
	formatters_by_ft = {
		go = { "goimports", "gofmt" },
		typescript = { "oxfmt" },
		typescriptreact = { "oxfmt" },
		javascript = { "oxfmt" },
		javascriptreact = { "oxfmt" },
		json = { "oxfmt" },
		lua = { "stylua" },
	},
	-- Format when the buffer is written
	format_on_save = { timeout_ms = 3000 },
})
