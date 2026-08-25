return {
  'saghen/blink.cmp',
  version = '1.*',
  event = { 'BufReadPre', 'BufNewFile' },
  dependencies = {
    { 'antosha417/nvim-lsp-file-operations', config = true },
    {
      'folke/lazydev.nvim',
      opts = {
        library = {
          -- load snacks types when `Snacks` is mentioned, for the global set in snacks/init.lua
          { path = 'snacks.nvim', words = { 'Snacks' } },
        },
      },
    },
  },
  ---@module 'blink.cmp'
  ---@type blink.cmp.Config
  opts = {
    keymap = {
      preset = 'default',
      ['<CR>'] = { 'accept', 'fallback' },
    },
    appearance = {
      -- 'mono' (default) for 'Nerd Font Mono' or 'normal' for 'Nerd Font'
      -- Adjusts spacing to ensure icons are aligned
      nerd_font_variant = 'mono',
    },
    fuzzy = { implementation = 'prefer_rust_with_warning' },
    signature = { enabled = true, window = { border = 'rounded' } },
    completion = {
      accept = { auto_brackets = { enabled = false } },
      documentation = {
        auto_show = true,
        auto_show_delay_ms = 250,
        window = {
          border = 'rounded',
        },
      },
      menu = {
        border = 'rounded',
        draw = {
          gap = 1,
          columns = { { 'kind_icon' }, { 'label', 'label_description', gap = 1 }, { 'kind' } },
        },
      },
      -- Display a preview of the selected item on the current line
      ghost_text = { enabled = true },
    },
    sources = {
      default = { 'lsp', 'path', 'buffer' },
      per_filetype = {
        codecompanion = { 'codecompanion' },
      },
    },
  },
}
