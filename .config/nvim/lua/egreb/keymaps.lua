vim.keymap.set('n', '<leader>cu', ':update<CR> :source<CR>')
vim.keymap.set('n', '<leader>w', ':write<CR>')
vim.keymap.set('n', '<leader>q', ':quit<CR>')

-- clear current search
vim.keymap.set('n', '<Esc>', '<cmd>nohlsearch<CR>')

-- Show error message
vim.keymap.set('n', '<leader>em', vim.diagnostic.open_float, { desc = 'Show diagnostic [E]rror [M]essages' })
-- Diagnostic keymaps
vim.keymap.set('n', '<leader>eq', vim.diagnostic.setloclist, { desc = 'Open diagnostic [E]rror [Q]uickfixlist' })

local map = vim.keymap.set

-- buffers
map('n', '<S-h>', '<cmd>bprevious<cr>', { desc = 'Prev Buffer' })
map('n', '<S-l>', '<cmd>bnext<cr>', { desc = 'Next Buffer' })
map('n', '[b', '<cmd>bprevious<cr>', { desc = 'Prev Buffer' })
map('n', ']b', '<cmd>bnext<cr>', { desc = 'Next Buffer' })
map('n', '<leader>bb', '<cmd>e #<cr>', { desc = 'Switch to Other Buffer' })
map('n', '<leader>`', '<cmd>e #<cr>', { desc = 'Switch to Other Buffer' })
map('n', '<leader>bd', function()
	-- @Snacks
	Snacks.bufdelete()
end, { desc = 'Delete Buffer' })
map('n', '<leader>bo', function()
	Snacks.bufdelete.other()
end, { desc = 'Delete Other Buffers' })
map('n', '<leader>bD', '<cmd>:bd<cr>', { desc = 'Delete Buffer and Window' })

-- multicursor keybindings
-- See plugin

-- Picker keymaps
-- find files
map({ 'n' }, '<C-p>', function()
	Snacks.picker.pick {
		source = 'files',
		-- cwd = vim.fn.getcwd(),
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[F]ind [F]iles' })
map({ 'n' }, '<leader>ff', function()
	Snacks.picker.pick {
		source = 'files',
		-- cwd = vim.fn.getcwd(),
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[F]ind [F]iles' })
-- grep files
map({ 'n' }, '<leader>fg', function()
	Snacks.picker.grep {
		cwd = true,
		hidden = true,
		ignored = false,
	}
end, { desc = '[G]rep [F]iles' })
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
		cwd = true,
		hidden = true,
		ignored = false,
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
