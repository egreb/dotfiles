vim.pack.add({
	"https://github.com/stevearc/oil.nvim",
})

-- File explorer that edits the filesystem like a buffer. Toggled with `<C-e>`, see keymaps.
require("oil").setup({
	default_file_explorer = true, -- replace netrw
	view_options = {
		show_hidden = true, -- show dotfiles and dot directories (toggle with `g.`)
	},
})
