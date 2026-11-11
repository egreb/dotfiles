vim.pack.add({
	"https://github.com/MeanderingProgrammer/render-markdown.nvim",
})

-- Prettier view of markdown files, toggled with `<leader>mo`, see keymaps.
-- Inspired by https://tduyng.com/blog/neovim-markdown-notes/. Everything not set here uses the plugin
-- defaults (full width heading backgrounds, GitHub style tables, Obsidian style callouts etc.).
require("render-markdown").setup({
	-- Completions for checkboxes and callouts through an in-process LSP (shows up in blink's `lsp` source)
	completions = { lsp = { enabled = true } },

	-- By default the rendering is only shown in normal/command/terminal mode and turns raw in insert mode.
	-- `true` keeps each element rendered in every mode, so the buffer doesn't jump while typing.
	-- (The line the cursor is on is still shown raw, that is anti conceal.)
	heading = { render_modes = true },
	pipe_table = { render_modes = true },
	checkbox = {
		render_modes = true,
		custom = {
			-- Personal todo marker: `[-]`
			todo = { raw = "[-]", rendered = "󰗡 Todo", highlight = "RenderMarkdownTodo" },
		},
	},
	code = {
		render_modes = true,
		border = "thin", -- thin line above/below code blocks instead of hiding the fences
	},
})
