-- plugins
vim.pack.add({ "https://github.com/folke/snacks.nvim" })

-- options
local Snacks = require("snacks")

Snacks.setup({
	dim = { enabled = true },
	indent = { enabled = true },
	scroll = { enabled = true },
	notifier = { enabled = true },
	picker = {
		enabled = true,
		hidden = false,
		ignored = false,
		cwd = vim.fn.getcwd(),
		layout = { preset = 'ivy', hidden = { 'preview' } },
		sources = {
			explorer = {
				layout = {
					layout = {
						position = 'right',
					},
				},
			},
		},
	},
	explorer = {
		enabled = true,
	},
	input = {
		enabled = true,
		icon = ' ',
		icon_hl = 'SnacksInputIcon',
		icon_pos = 'left',
		prompt_pos = 'title',
		win = { style = 'input' },
		expand = true,
	},
	scratch = { enabled = true },
	lazygit = {
		enabled = true,
		configure = true,
		win = {
			style = 'lazygit',
		},
		selectedLineBgColor = { bg = '#333' },
	},
	gitbrowse = { enabled = true },
	bufdelete = {},
})

-- keymaps
local keys = {
	{
		'<C-n>',
		function()
			Snacks.explorer()
		end,
		desc = 'Explorer',
	},
	{
		'<leader>nn',
		function()
			Snacks.scratch()
		end,
		desc = 'Toggle Scratch Buffer',
	},
	{
		'<leader>ns',
		function()
			Snacks.scratch.select()
		end,
		desc = 'Select Scratch Buffer',
	},
	{
		'<leader>lg',
		function()
			Snacks.lazygit.open()
		end,
		desc = '[L]azy[G]it',
	},
	{
		'<leader>gb',
		function()
			Snacks.git.blame_line()
		end,
		desc = '[G]it [B]lame',
	},
	{
		'<leader>go',
		function()
			Snacks.gitbrowse.open()
		end,
		desc = '[G]it [O]pen',
	},
	{
		'<leader>tt',
		function()
			Snacks.terminal.toggle()
		end,
		desc = '[T]oggle [T]erminal',
	},
}

for _, key in ipairs(keys) do
	vim.keymap.set("n", key[1], key[2], { desc = key.desc })
end

local map = vim.keymap.set

map('n', '<S-h>', '<cmd>bprevious<cr>', { desc = 'Prev Buffer' })
map('n', '<S-l>', '<cmd>bnext<cr>', { desc = 'Next Buffer' })
map('n', '[b', '<cmd>bprevious<cr>', { desc = 'Prev Buffer' })
map('n', ']b', '<cmd>bnext<cr>', { desc = 'Next Buffer' })
map('n', '<leader>bb', '<cmd>e #<cr>', { desc = 'Switch to Other Buffer' })
map('n', '<leader>`', '<cmd>e #<cr>', { desc = 'Switch to Other Buffer' })
map('n', '<leader>bd', function()
	Snacks.bufdelete()
end, { desc = 'Delete Buffer' })
map('n', '<leader>bo', function()
	Snacks.bufdelete.other()
end, { desc = 'Delete Other Buffers' })
map('n', '<leader>bD', '<cmd>:bd<cr>', { desc = 'Delete Buffer and Window' })

-- find files
map({ 'n' }, '<C-p>', function()
	Snacks.picker.pick {
		source = 'files',
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[F]ind [F]iles' })
map({ 'n' }, '<leader>ff', function()
	Snacks.picker.pick {
		source = 'files',
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[F]ind [F]iles' })
-- grep files
map({ 'n' }, '<leader>fG', function()
	Snacks.picker.grep {
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[G]rep [F]iles' })
-- grep files, literal (no regex) — for paths/patterns with {}, (), [], etc.
map({ 'n' }, '<leader>fg', function()
	Snacks.picker.grep {
		cwd = true,
		hidden = true,
		ignored = false,
		regex = false,
	}
end, { desc = '[G]rep [F]iles (literal)' })
-- resume previous search
map({ 'n' }, '<leader>fr', function()
	Snacks.picker.resume { cwd = true }
end, { desc = '[R]esume [S]earch' })

-- list buffers
map({ 'n' }, '<leader><leader>', function()
	Snacks.picker.buffers { cwd = true }
end, { desc = '[F]ind [B]uffers' })

-- search history
map({ 'n' }, '<leader>f.', function()
	Snacks.picker.recent {
		cwd = true,
		hidden = true,
		ignored = false,
		filter = { cwd = true },
	}
end, { desc = 'Search [H]istory' })
-- diagnostics
map({ 'n' }, '<leader>fd', function()
	Snacks.picker.diagnostics {
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = 'Find [D]iagnostic' })
-- old files
map({ 'n' }, '<leader>fo', function()
	Snacks.picker.recent {
		hidden = true,
		ignored = false,
		filter = { cwd = true },
	}
end, { desc = '[O]ld Files' })
-- search word under cursor
map({ 'n' }, '<leader>fw', function()
	Snacks.picker.grep_word {
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = 'Visual selection or word' })
map({ 'n' }, '<leader>fb', function()
	Snacks.picker.lines { pattern = vim.fn.expand '<cword>' }
end, { desc = '[F]ind in [B]uffer' })
map({ 'n' }, '<leader>fs', function()
	Snacks.picker.lsp_symbols {}
end, { desc = '[F]ind [S]ymbols' })
