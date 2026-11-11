local M = {}

--- Yank the path of the current file with the line number (or line range when in visual mode).
---@param absolute boolean absolute path instead of relative to the cwd
function M.path_with_line(absolute)
	local mode = vim.fn.mode()
	local first, last = vim.fn.line("."), vim.fn.line(".")
	if mode == "v" or mode == "V" or mode == "\22" then
		first, last = vim.fn.line("v"), vim.fn.line(".")
		if first > last then
			first, last = last, first
		end
		-- leave visual mode
		vim.api.nvim_feedkeys(vim.keycode("<Esc>"), "n", false)
	end

	local path = vim.fn.expand(absolute and "%:p" or "%:.")
	local text = first == last and ("%s:%d"):format(path, first) or ("%s:%d-%d"):format(path, first, last)

	vim.fn.setreg("+", text)
	vim.notify("Yanked " .. text)
end

return M
