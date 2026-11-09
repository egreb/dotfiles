local augroup = vim.api.nvim_create_augroup("egreb", { clear = true })

-- Highlight yanked lines
vim.api.nvim_create_autocmd("TextYankPost", {
	group = augroup,
	desc = "Highlight yanked text",
	callback = function()
		vim.hl.on_yank()
	end,
})

-- Disable automatic comment on newline
vim.api.nvim_create_autocmd("FileType", {
	group = augroup,
	desc = "Disable automatic comment on newline",
	callback = function()
		vim.opt_local.formatoptions:remove({ "c", "r", "o" })
	end,
})

-- Restore file cursor to position in previous session
vim.api.nvim_create_autocmd("BufReadPost", {
	group = augroup,
	desc = "Restore cursor position from previous session",
	callback = function(args)
		local mark = vim.api.nvim_buf_get_mark(args.buf, '"')
		local line_count = vim.api.nvim_buf_line_count(args.buf)
		if mark[1] > 0 and mark[1] <= line_count then
			vim.api.nvim_win_set_cursor(0, mark)
		end
	end,
})
