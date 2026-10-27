-- plugins
vim.pack.add({
	{ src = "https://github.com/jake-stewart/multicursor.nvim", version = "1.0" },
})

-- options
local mc = require("multicursor-nvim")

mc.setup()

-- keymaps
local set = vim.keymap.set

-- Add and remove cursors with control + left click.
set("n", "<c-leftmouse>", mc.handleMouse)
set("n", "<c-leftdrag>", mc.handleMouseDrag)
set("n", "<c-leftrelease>", mc.handleMouseRelease)

set("x", "I", mc.insertVisual)
set({ "n", "v" }, "[d", function()
	mc.lineAddCursor(-1)
end, { desc = "Add cursor above" })
set({ "n", "v" }, "]d", function()
	mc.lineAddCursor(1)
end, { desc = "Add cursor below" })

-- Add a new cursor by matching word/selection
set({ "n", "v" }, "[m", function()
	mc.matchAddCursor(1)
end, { desc = "Add cursor previous match" })
set({ "n", "v" }, "]m", function()
	mc.matchAddCursor(-1)
end, { desc = "Add cursor next match" })

set("n", "<c-q>", function()
	if not mc.cursorsEnabled() then
		mc.enableCursors()
	elseif mc.hasCursors() then
		mc.clearCursors()
	end
end, { desc = "Enable/clear cursors" })
set({ "n", "x" }, "ga", mc.addCursorOperator, { desc = "Add cursors over motion" })

-- Customize how cursors look.
local hl = vim.api.nvim_set_hl
hl(0, "MultiCursorCursor", { reverse = true })
hl(0, "MultiCursorVisual", { link = "Visual" })
hl(0, "MultiCursorSign", { link = "SignColumn" })
hl(0, "MultiCursorMatchPreview", { link = "Search" })
hl(0, "MultiCursorDisabledCursor", { reverse = true })
hl(0, "MultiCursorDisabledVisual", { link = "Visual" })
hl(0, "MultiCursorDisabledSign", { link = "SignColumn" })
