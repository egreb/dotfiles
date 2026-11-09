local map = vim.keymap.set

-- Git
map("n", "]g", function()
	require("gitsigns").nav_hunk("next")
end, { desc = "Next git hunk" })

map("n", "[g", function()
	require("gitsigns").nav_hunk("prev")
end, { desc = "Previous git hunk" })

map("n", "<leader>cd", function()
	require("egreb.utils.codediff").toggle()
end, { desc = "Toggle codediff" })

-- Opens the current line (or the visual selection) on the git remote, e.g. GitHub
map({ "n", "v" }, "<leader>go", function()
	require("snacks").gitbrowse({ what = "file" })
end, { desc = "Open line on GitHub" })
