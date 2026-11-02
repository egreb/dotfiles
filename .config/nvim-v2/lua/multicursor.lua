-- helpers for the native multicursor (:help multicursor)
-- extra cursors are extmarks in this namespace; Q, gQ, ]C, [C and zq stay native
local ns = vim.api.nvim_create_namespace("nvim.multicursor")

-- mappings are replayed at every cursor, so only act for the primary one
local function primary(fn)
	return function()
		if not vim.v.cascading then
			fn()
		end
	end
end

local function has_cursors()
	return #vim.api.nvim_buf_get_extmarks(0, ns, 0, -1, { limit = 1 }) > 0
end

-- search for the visual selection (or the word under the cursor) and put the cursor on its first match
local function set_pattern()
	local mode = vim.fn.mode()
	local pattern
	if mode == "v" or mode == "V" or mode == "\22" then
		local from, to = vim.fn.getpos("v"), vim.fn.getpos(".")
		local text = table.concat(vim.fn.getregion(from, to, { type = mode }), "\n")
		pattern = "\\C\\V" .. vim.fn.escape(text, "\\"):gsub("\n", "\\n")
		if from[2] > to[2] or (from[2] == to[2] and from[3] > to[3]) then
			from = to
		end
		vim.cmd.normal({ "\27", bang = true })
		vim.api.nvim_win_set_cursor(0, { from[2], from[3] - 1 })
	else
		local word = vim.fn.expand("<cword>")
		if word == "" then
			return nil
		end
		pattern = "\\C\\V\\<" .. vim.fn.escape(word, "\\") .. "\\>"
		local line = vim.fn.line(".")
		if vim.fn.search(pattern, "bc", line) == 0 then
			vim.fn.search(pattern, "c", line)
		end
	end
	vim.fn.setreg("/", pattern)
	vim.v.hlsearch = 1
	return pattern
end

-- jump to the next match, leaving a cursor behind unless skipping
local function next_match(add)
	local pattern = (has_cursors() and vim.fn.mode() == "n") and vim.fn.getreg("/") or set_pattern()
	if not pattern or pattern == "" then
		return
	end
	local pos = vim.api.nvim_win_get_cursor(0)
	if add then
		vim.api.nvim_buf_set_extmark(0, ns, pos[1] - 1, pos[2], {})
	end
	if vim.fn.search(pattern, "w") == 0 then
		vim.notify("No matches", vim.log.levels.WARN)
		return
	end
	local row, col = unpack(vim.api.nvim_win_get_cursor(0))
	if #vim.api.nvim_buf_get_extmarks(0, ns, { row - 1, col }, { row - 1, col }, {}) > 0 then
		vim.notify("No more matches", vim.log.levels.INFO)
	end
end

local function select_all()
	local pattern = set_pattern()
	if not pattern then
		return
	end
	local row, col = unpack(vim.api.nvim_win_get_cursor(0))
	for _, m in ipairs(vim.fn.matchbufline(0, pattern, 1, "$")) do
		if m.lnum ~= row or m.byteidx ~= col then
			vim.api.nvim_buf_set_extmark(0, ns, m.lnum - 1, m.byteidx, {})
		end
	end
end

vim.keymap.set({ "n", "x" }, "<A-n>", primary(function()
	next_match(true)
end), { desc = "Multicursor: select next occurrence" })
vim.keymap.set({ "n", "x" }, "<leader>ms", primary(function()
	next_match(false)
end), { desc = "Multicursor: skip next occurrence" })
vim.keymap.set({ "n", "x" }, "<leader>ma", primary(select_all), { desc = "Multicursor: select all matches" })
vim.keymap.set("x", "<leader>ml", "Q", { desc = "Multicursor: add cursors to selected lines" })

-- <Esc> (and <C-[>, the same key) clears search highlight and cursors
vim.keymap.set("n", "<Esc>", primary(function()
	vim.cmd.nohlsearch()
	vim.api.nvim_buf_clear_namespace(0, ns, 0, -1)
end), { desc = "Clear search highlight and multicursors" })
