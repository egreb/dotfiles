vim.pack.add({
	{
		src = "https://github.com/obsidian-nvim/obsidian.nvim",
		version = vim.version.range("*"), -- latest release
	},
})

require("obsidian").setup({
	legacy_commands = false, -- use the `:Obsidian <subcommand>` commands only
	-- The first workspace is the default one
	workspaces = {
		{ name = "work", path = "~/vaults/work" },
		{ name = "personal", path = "~/vaults/personal" },
	},
	picker = { name = "snacks.picker" },
	-- md-render handles the rendered markdown view (plugins/md-render.lua)
	ui = { enable = false },
})
