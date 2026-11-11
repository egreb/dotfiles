vim.pack.add({
	"https://github.com/zhisme/copy_with_context.nvim",
})

-- Copies the line (or visual selection) together with where it came from. The plugin always
-- includes the code, so the path-only yanks (`<leader>yf`, `<leader>ya`) live in utils/yank.lua.
require("copy_with_context").setup({
	mappings = {
		relative = "<leader>yy", -- code + relative path and line(s)
		absolute = "<leader>yY", -- code + absolute path and line(s)
	},
	output_formats = {
		-- used by both mappings above
		default = "{copied_text}\n\n{filepath}:{line}",
	},
})
