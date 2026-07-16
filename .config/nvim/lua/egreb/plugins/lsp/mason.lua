return {
  {
    'williamboman/mason-lspconfig.nvim',
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
    opts = {
      ensure_installed = {
        'prettier', -- prettier formatter
        'stylua', -- lua formatter
        'goimports',
        'djlint', -- go template (html) formatter
      },
    },
    dependencies = {
      'williamboman/mason.nvim',
    },
  },
}
