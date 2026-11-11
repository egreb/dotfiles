vim.pack.add({
	"https://github.com/lewis6991/gitsigns.nvim",
	"https://github.com/esmuellert/codediff.nvim",
})

-- Git integration for the buffer (signs in the sign column, hunk navigation)
require("gitsigns").setup({})

-- Review local diffs / do code reviews. Opened with `:CodeDiff`, see `<leader>cd`.
-- The diff library is downloaded automatically on first use.
