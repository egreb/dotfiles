return {
	'obsidian-nvim/obsidian.nvim',
	version = '*',
	lazy = true,
	ft = 'markdown',
	dependencies = {
		'nvim-lua/plenary.nvim',
	},
	opts = {
		daily_notes = {
			-- Optional, if you keep daily notes in a separate directory.
			folder = 'notes/dailies',
			-- Optional, if you want to change the date format for the ID of daily notes.
			date_format = '%Y-%m-%d',
			-- Optional, if you want to change the date format of the default alias of daily notes.
			alias_format = '%B %-d, %Y',
			-- Optional, default tags to add to each new daily note created.
			default_tags = { 'daily-notes' },
			-- Optional, if you want to automatically insert a template from your template directory like 'daily.md'
			template = nil,
		},
		workspaces = {
			{
				name = 'personal',
				path = '~/vaults/personal',
			},
			{
				name = 'work',
				path = '~/vaults/work',
			},
		},
	},
}
