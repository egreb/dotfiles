-- options (set before the plugin loads so dadbod-ui picks them up)
vim.g.db_ui_use_nerd_fonts = 1

-- plugins
vim.pack.add({
	"https://github.com/tpope/vim-dadbod",
	"https://github.com/kristijanhusak/vim-dadbod-ui",
	"https://github.com/kristijanhusak/vim-dadbod-completion",
})

-- keymaps
vim.keymap.set("n", "<leader>db", "<cmd>DBUIToggle<CR>", { desc = "Toggle database UI" })
