vim.pack.add({
	"https://github.com/nvim-mini/mini.nvim",
})

require("mini.pairs").setup() -- auto close brackets and quotes
require("mini.surround").setup() -- add/delete/replace surroundings (sa, sd, sr)
require("mini.ai").setup() -- extended a/i text objects
require("mini.diff").setup() -- visualize diff hunks (colors the line numbers)
