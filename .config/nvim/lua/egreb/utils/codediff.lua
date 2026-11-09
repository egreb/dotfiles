local M = {}

-- Tabpages that currently show a CodeDiff view
local tabs = {}

vim.api.nvim_create_autocmd("User", {
	group = vim.api.nvim_create_augroup("egreb_codediff", { clear = true }),
	pattern = { "CodeDiffOpen", "CodeDiffClose" },
	desc = "Track which tabpages are showing CodeDiff",
	callback = function(event)
		tabs[event.data.tabpage] = event.match == "CodeDiffOpen" or nil
	end,
})

--- Close CodeDiff if the current tabpage is showing it, otherwise open it.
function M.toggle()
	if tabs[vim.api.nvim_get_current_tabpage()] then
		vim.cmd.tabclose()
	else
		vim.cmd("CodeDiff")
	end
end

return M
