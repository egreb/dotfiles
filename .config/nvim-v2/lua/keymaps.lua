-- jj to replace escape
vim.keymap.set("i", "jj", "<ESC>", { silent = true })
vim.keymap.set("n", "<leader>cu", ":update<CR> :source<CR>")
vim.keymap.set("n", "<leader>w", ":write<CR>")
vim.keymap.set("n", "<leader>q", ":quit<CR>")
-- clear current search
vim.keymap.set("n", "<Esc>", "<cmd>nohlsearch<CR>")
