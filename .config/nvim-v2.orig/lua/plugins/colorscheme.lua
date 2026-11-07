-- Catppuccin is bundled with Neovim; dark selects Mocha.
vim.opt.termguicolors = true
vim.opt.background = "dark"
vim.cmd.colorscheme("catppuccin")

-- gopls string tokens have higher priority than Tree-sitter SQL injections.
-- Let Tree-sitter color Go strings, including any embedded languages.
local function clear_go_string_semantic_highlight()
	vim.api.nvim_set_hl(0, "@lsp.type.string.go", {})
end

vim.api.nvim_create_autocmd("ColorScheme", {
	group = vim.api.nvim_create_augroup("go_string_highlights", { clear = true }),
	callback = clear_go_string_semantic_highlight,
})
clear_go_string_semantic_highlight()
