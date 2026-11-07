-- plugins
vim.pack.add({
	"https://github.com/lewis6991/gitsigns.nvim",
	"https://github.com/martindur/zdiff.nvim",
})

-- options
require("gitsigns").setup({
	current_line_blame = true,
})

require("zdiff").setup()

-- keymaps
vim.keymap.set("n", "[h", function()
	require("gitsigns").nav_hunk("prev")
end, { desc = "Previous Git hunk" })
vim.keymap.set("n", "]h", function()
	require("gitsigns").nav_hunk("next")
end, { desc = "Next Git hunk" })

-- open the GitHub PR that introduced the current line (or the sha under the cursor in :Gitsigns blame)
vim.keymap.set("n", "<leader>gp", function()
	local sha = vim.fn.expand("<cword>")
	if not sha:match("^%x%x%x%x%x%x%x+$") then
		local file = vim.api.nvim_buf_get_name(0)
		local out = vim.system({ "git", "blame", "-l", "-L", vim.fn.line(".") .. ",+1", "--", file }, { cwd = vim.fs.dirname(file) }):wait()
		sha = out.code == 0 and out.stdout:match("^%^?(%x+)") or nil
		if not sha or sha:match("^0+$") then
			vim.notify("No commit for this line", vim.log.levels.WARN)
			return
		end
	end
	vim.system({ "gh", "pr", "list", "--state", "all", "--search", sha, "--json", "url", "--jq", ".[0].url" }, { text = true }, function(out)
		local url = vim.trim(out.stdout or "")
		vim.schedule(function()
			if out.code ~= 0 then
				vim.notify("gh: " .. vim.trim(out.stderr or ""), vim.log.levels.ERROR)
			elseif url == "" then
				vim.notify("No PR found for " .. sha:sub(1, 8), vim.log.levels.WARN)
			else
				vim.ui.open(url)
			end
		end)
	end)
end, { desc = "Open PR for line's commit" })

vim.keymap.set("n", "<leader>gm",":Gitsigns diffthis main<cr>", { silent = true, desc = "Diff against main" })

vim.keymap.set("n", "<leader>zd", function()
	require("zdiff").open()
end, { desc = "Zdiff (uncommitted)" })
vim.keymap.set("n", "<leader>zD", function()
	require("zdiff").open("main")
end, { desc = "Zdiff (vs main)" })
