vim.pack.add({
	{
		src = "https://github.com/delphinus/md-render.nvim",
		version = vim.version.range("*"), -- latest release
	},
})

-- Rendered markdown (callouts, tables, code blocks, images, wikilinks, ...). Unlike an in-buffer renderer
-- it never touches the file's buffer: "read mode" swaps the window over to a separate read-only render
-- buffer, and toggling again swaps back to the source for editing.
-- Markdown files open as plain source, ready for editing; read mode is opt-in with `<leader>mo` (see keymaps).
-- No setup call needed, the plugin works with its defaults.
