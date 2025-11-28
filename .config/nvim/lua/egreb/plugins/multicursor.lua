return {
	'jake-stewart/multicursor.nvim',
	branch = '1.0',
	config = function()
		local mc = require 'multicursor-nvim'
		--
		mc.setup()
		local set = vim.keymap.set
		--
		-- -- Add and remove cursors with control + left click.
		set('n', '<c-leftmouse>', mc.handleMouse)
		set('n', '<c-leftdrag>', mc.handleMouseDrag)
		set('n', '<c-leftrelease>', mc.handleMouseRelease)

		set('x', 'I', mc.insertVisual)
		set({ 'n', 'v' }, '[n', function()
			mc.lineAddCursor(-1)
		end)
		set({ 'n', 'v' }, ']n', function()
			mc.lineAddCursor(1)
		end)
		--
		-- Add or skip adding a new cursor by matching word/selection
		set({ 'n', 'v' }, '[N', function()
			mc.matchAddCursor(1)
		end, {
			desc = 'Add Cursor [P]revious match',
		})
		set({ 'n', 'v' }, ']N', function()
			mc.matchAddCursor(-1)
		end, {
			desc = 'Add Cursor [N]ext match',
		})
		-- set({ 'n', 'v' }, '<leader>ls', function()
		--   mc.matchSkipCursor(1)
		-- end)
		-- set({ 'n', 'v' }, '<leader>lS', function()
		--   mc.matchSkipCursor(-1)
		-- end)

		set('n', '<c-q>', function()
			if not mc.cursorsEnabled() then
				mc.enableCursors()
			elseif mc.hasCursors() then
				mc.clearCursors()
			else
				-- Default <esc> handler.
			end
		end)

		-- Customize how cursors look.
		local hl = vim.api.nvim_set_hl
		hl(0, 'MultiCursorCursor', { reverse = true })
		hl(0, 'MultiCursorVisual', { link = 'Visual' })
		hl(0, 'MultiCursorSign', { link = 'SignColumn' })
		hl(0, 'MultiCursorMatchPreview', { link = 'Search' })
		hl(0, 'MultiCursorDisabledCursor', { reverse = true })
		hl(0, 'MultiCursorDisabledVisual', { link = 'Visual' })
		hl(0, 'MultiCursorDisabledSign', { link = 'SignColumn' })
	end,
}
