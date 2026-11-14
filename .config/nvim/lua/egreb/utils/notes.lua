local M = {}

-- Root of the active Obsidian workspace, falling back to the work vault
local function vault_root()
	local ws = rawget(_G, "Obsidian") and Obsidian.workspace
	if ws and ws.root then
		return tostring(ws.root)
	end
	return vim.fn.expand("~/vaults/work")
end

--- Open the most recently modified note in the vault.
function M.open_recent()
	local root = vim.fs.normalize(vault_root())
	local newest, newest_mtime = nil, -1
	local notes = vim.fs.find(function(name)
		return name:match("%.md$") ~= nil
	end, { path = root, type = "file", limit = math.huge })
	for _, path in ipairs(notes) do
		-- skip hidden folders (.obsidian, .trash, ...)
		if not path:sub(#root + 1):find("/%.") then
			local stat = vim.uv.fs_stat(path)
			local mtime = stat and (stat.mtime.sec + stat.mtime.nsec / 1e9) or -1
			if mtime > newest_mtime then
				newest, newest_mtime = path, mtime
			end
		end
	end
	if not newest then
		vim.notify("No notes found in " .. root, vim.log.levels.WARN)
		return
	end
	vim.cmd.edit(vim.fn.fnameescape(newest))
end

return M
