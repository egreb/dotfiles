local map = vim.keymap.set

-- Basics
map("n", "<leader>w", "<cmd>write<cr>", { desc = "Write buffer" })
map("n", "<leader>q", "<cmd>bdelete<cr>", { desc = "Close buffer" })

-- Navigate between open windows
map("n", "<C-h>", "<C-w>h", { desc = "Window left" })
map("n", "<C-j>", "<C-w>j", { desc = "Window down" })
map("n", "<C-k>", "<C-w>k", { desc = "Window up" })
map("n", "<C-l>", "<C-w>l", { desc = "Window right" })

-- Move line(s) up/down (opt+j/k, needs the terminal to send option as alt)
map("n", "<A-j>", "<cmd>move .+1<cr>==", { desc = "Move line down" })
map("n", "<A-k>", "<cmd>move .-2<cr>==", { desc = "Move line up" })
map("v", "<A-j>", ":move '>+1<cr>gv=gv", { desc = "Move selection down", silent = true })
map("v", "<A-k>", ":move '<-2<cr>gv=gv", { desc = "Move selection up", silent = true })

-- Yank path with line number. `<leader>yy` (code with context) comes from copy_with_context.
map({ "n", "v" }, "<leader>yf", function()
	require("egreb.utils.yank").path_with_line(false)
end, { desc = "Yank relative path with line number" })

map({ "n", "v" }, "<leader>ya", function()
	require("egreb.utils.yank").path_with_line(true)
end, { desc = "Yank absolute path with line number" })

-- LSP
map("n", "<leader>ih", function()
	vim.lsp.inlay_hint.enable(not vim.lsp.inlay_hint.is_enabled())
end, { desc = "Toggle inlay hints" })

-- LSP navigation (Snacks pickers: jump straight to the target if there is only one, otherwise pick with preview)
map("n", "gd", function()
	require("snacks").picker.lsp_definitions()
end, { desc = "Goto definition" })

map("n", "gi", function()
	require("snacks").picker.lsp_implementations()
end, { desc = "Goto implementation" })

-- Overrides the builtin `grr` (quickfix list) with the Snacks picker
map("n", "grr", function()
	require("snacks").picker.lsp_references()
end, { desc = "Goto references" })

-- Find (Snacks picker, shows a preview of the selected item)
map("n", "<leader>ff", function()
	require("snacks").picker.files()
end, { desc = "Find files" })

map("n", "<leader>fg", function()
	require("snacks").picker.grep()
end, { desc = "Grep files" })

map("n", "<leader>f.", function()
	require("snacks").picker.recent()
end, { desc = "Search history" })

map("n", "<leader>fr", function()
	require("snacks").picker.resume()
end, { desc = "Resume search" })

map("n", "<leader><leader>", function()
	require("snacks").picker.buffers()
end, { desc = "List buffers" })

map("n", "<leader>fd", function()
	require("snacks").picker.diagnostics()
end, { desc = "LSP diagnostics" })

map("n", "<leader>fs", function()
	require("snacks").picker.lsp_symbols()
end, { desc = "List functions and symbols" })

-- File explorer
map("n", "<C-e>", function()
	require("oil").toggle_float()
end, { desc = "Toggle Oil explorer" })

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

-- Notes (Obsidian, defaults to the work vault)
map("n", "<leader>no", "<cmd>Obsidian quick_switch<cr>", { desc = "Open note" })
map("n", "<leader>ns", "<cmd>Obsidian search<cr>", { desc = "Search notes" })
map("n", "<leader>nn", "<cmd>Obsidian new<cr>", { desc = "New note" })

-- Markdown
map("n", "<leader>mo", "<cmd>RenderMarkdown toggle<cr>", { desc = "Toggle markdown rendering" })

map("n", "<leader>ml", function()
	vim.snippet.expand("[${1:text}](${2:url})$0")
end, { desc = "Insert markdown link" })

map("n", "<leader>mt", function()
	vim.snippet.expand("| ${1:Header} | ${2:Header} |\n| --- | --- |\n| ${3} | ${4} |$0")
end, { desc = "Insert markdown table" })
