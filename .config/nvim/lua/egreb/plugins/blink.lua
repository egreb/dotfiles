vim.pack.add({
	-- Pinned to the 1.x release tags: on a release tag blink downloads a prebuilt fuzzy matcher,
	-- so no Rust toolchain is needed (tracking `main` requires building from source).
	{ src = "https://github.com/saghen/blink.cmp", version = vim.version.range("1.*") },
})

require("blink.cmp").setup({
	-- 'default': C-y accepts, C-n/C-p or Up/Down select, C-space opens menu/docs, C-e hides
	keymap = { preset = "default" },
	completion = { documentation = { auto_show = false } }, -- docs only on demand
	sources = { default = { "lsp", "path", "snippets", "buffer" } },
	fuzzy = { implementation = "prefer_rust" }, -- prebuilt binary, falls back to lua if unavailable
})

-- Tell every language server that the client supports blink's completion features
vim.lsp.config("*", {
	capabilities = require("blink.cmp").get_lsp_capabilities(),
})
