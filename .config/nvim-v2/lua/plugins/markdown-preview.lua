-- The plugin ships a node app; fetch its prebuilt binary after install/update.
-- Must be registered before vim.pack.add so it fires on first install.
vim.api.nvim_create_autocmd("PackChanged", {
	callback = function(ev)
		local name, kind = ev.data.spec.name, ev.data.kind
		if name == "markdown-preview.nvim" and (kind == "install" or kind == "update") then
			if not ev.data.active then
				vim.cmd.packadd(name)
			end
			vim.fn["mkdp#util#install"]()
		end
	end,
})

vim.pack.add({ "https://github.com/iamcco/markdown-preview.nvim" })

-- Keep the preview open when switching to another buffer.
vim.g.mkdp_auto_close = 0

vim.keymap.set("n", "<leader>mp", "<cmd>MarkdownPreviewToggle<CR>", { desc = "Toggle Markdown preview" })
