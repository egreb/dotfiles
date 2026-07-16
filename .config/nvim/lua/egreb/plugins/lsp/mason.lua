return {
  {
    'williamboman/mason-lspconfig.nvim',
    event = { 'BufReadPre', 'BufNewFile' },
    cmd = 'Mason',
    opts = {
      -- list of servers for mason to install
      ensure_installed = {
        'tsgo',
        'html',
        'tailwindcss',
        'lua_ls',
        'emmet_ls',
        'jsonls',
        'bashls',
        'vue_ls',
        'gopls',
      },
    },
    dependencies = {
      {
        'williamboman/mason.nvim',
        opts = {
          ui = {
            icons = {
              package_installed = '✓',
              package_pending = '➜',
              package_uninstalled = '✗',
            },
          },
        },
      },
      {
        'neovim/nvim-lspconfig',
      },
    },
  },
  {
    'WhoIsSethDaniel/mason-tool-installer.nvim',
    event = 'VeryLazy',
    opts = {
      ensure_installed = {
        'prettier', -- prettier formatter
        'stylua', -- lua formatter
        'goimports',
        'djlint', -- go template (html) formatter
      },
    },
    config = function(_, opts)
      local mti = require 'mason-tool-installer'
      mti.setup(opts)
      -- The plugin normally runs its install check from a VimEnter autocmd;
      -- VeryLazy fires after VimEnter, so trigger it explicitly.
      mti.run_on_start()
    end,
    dependencies = {
      'williamboman/mason.nvim',
    },
  },
}
