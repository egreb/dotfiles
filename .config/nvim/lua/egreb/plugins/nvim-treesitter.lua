return {

	-- Treesitter is a new parser generator tool that we can
	-- use in Neovim to power faster and more accurate
	-- syntax highlighting.
	{
		'nvim-treesitter/nvim-treesitter',
		branch = 'main',
		version = false, -- last release is way too old and doesn't work on Windows
		build = function()
			local TS = require 'nvim-treesitter'
			if not TS.get_installed then
				print 'Please restart Neovim and run `:TSUpdate` to use the `nvim-treesitter` **main** branch.'
				return
			end
			print(TS.update(nil, { summary = true }))
		end,
		lazy = vim.fn.argc(-1) == 0, -- load treesitter early when opening a file from the cmdline
		event = { 'VeryLazy' },
		cmd = { 'TSUpdate', 'TSInstall', 'TSLog', 'TSUninstall' },
		opts_extend = { 'ensure_installed' },
		opts = {
			-- LazyVim config for treesitter
			indent = { enable = true },
			highlight = { enable = true },
			folds = { enable = true },
			ensure_installed = {
				'bash',
				'diff',
				'html',
				'javascript',
				'jsdoc',
				'json',
				'jsonc',
				'lua',
				'luadoc',
				'luap',
				'markdown',
				'markdown_inline',
				'printf',
				'query',
				'regex',
				'toml',
				'tsx',
				'typescript',
				'vim',
				'vimdoc',
				'xml',
				'yaml',
				'go',
			},
		},
		config = function(_, opts)
			local TS = require 'nvim-treesitter'

			setmetatable(require 'nvim-treesitter.install', {
				__newindex = function(_, k)
					if k == 'compilers' then
						vim.schedule(function()
							print {
								'Setting custom compilers for `nvim-treesitter` is no longer supported.',
								'',
								'For more info, see:',
								'- [compilers](https://docs.rs/cc/latest/cc/#compile-time-requirements)',
							}
						end)
					end
				end,
			})

			-- some quick sanity checks
			if not TS.get_installed then
				return print 'Please use `:Lazy` and update `nvim-treesitter`'
			elseif type(opts.ensure_installed) ~= 'table' then
				return print '`nvim-treesitter` opts.ensure_installed must be a table'
			end

			-- setup treesitter
			TS.setup(opts)

			-- install missing parsers
			local install = opts.ensure_installed or {}
			if #install > 0 then
				TS.install(install, { summary = true })
			end
		end,
	},

	{
		'nvim-treesitter/nvim-treesitter-textobjects',
		branch = 'main',
		event = 'VeryLazy',
		opts = {
			move = {
				enable = true,
				set_jumps = true, -- whether to set jumps in the jumplist
				-- LazyVim extention to create buffer-local keymaps
				keys = {
					goto_next_start = { [']f'] = '@function.outer', [']c'] = '@class.outer', [']a'] = '@parameter.inner' },
					goto_next_end = { [']F'] = '@function.outer', [']C'] = '@class.outer', [']A'] = '@parameter.inner' },
					goto_previous_start = { ['[f'] = '@function.outer', ['[c'] = '@class.outer', ['[a'] = '@parameter.inner' },
					goto_previous_end = { ['[F'] = '@function.outer', ['[C'] = '@class.outer', ['[A'] = '@parameter.inner' },
				},
			},
		},
		config = function(_, opts)
			local TS = require 'nvim-treesitter-textobjects'
			if not TS.setup then
				print 'Please use `:Lazy` and update `nvim-treesitter`'
				return
			end
			TS.setup(opts)

			vim.api.nvim_create_autocmd('FileType', {
				group = vim.api.nvim_create_augroup('egreb_treesitter_textobjects', { clear = true }),
				callback = function(ev)
					---@type table<string, table<string, string>>
					local moves = vim.tbl_get(opts, 'move', 'keys') or {}

					for method, keymaps in pairs(moves) do
						for key, query in pairs(keymaps) do
							local desc = query:gsub('@', ''):gsub('%..*', '')
							desc = desc:sub(1, 1):upper() .. desc:sub(2)
							desc = (key:sub(1, 1) == '[' and 'Prev ' or 'Next ') .. desc
							desc = desc .. (key:sub(2, 2) == key:sub(2, 2):upper() and ' End' or ' Start')
							if not (vim.wo.diff and key:find '[cC]') then
								vim.keymap.set({ 'n', 'x', 'o' }, key, function()
									require('nvim-treesitter-textobjects.move')[method](query, 'textobjects')
								end, {
									buffer = ev.buf,
									desc = desc,
									silent = true,
								})
							end
						end
					end
				end,
			})
		end,
	},
}
