return {
	cmd = { "lua-language-server" },
	filetypes = { "lua" },
	root_markers = { ".luarc.json", ".luarc.jsonc", ".stylua.toml", "stylua.toml", ".git" },
	settings = {
		Lua = {
			runtime = { version = "LuaJIT" }, -- neovim embeds LuaJIT
			workspace = {
				checkThirdParty = false,
				library = { vim.env.VIMRUNTIME }, -- know about `vim.*` so it isn't an undefined global
			},
		},
	},
}
