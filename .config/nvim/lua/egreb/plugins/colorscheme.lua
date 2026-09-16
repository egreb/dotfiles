return {
  'rose-pine/neovim',
  name = 'rose-pine',
  config = function()
    require('rose-pine').setup {
      -- 'auto' follows vim.o.background: dawn when light, dark_variant when dark.
      -- See the appearance autocmd in egreb.autocommands for what flips it.
      variant = 'auto',
      dark_variant = 'moon',
      styles = {
        bold = false,
        italic = false,
        transparency = false,
      },
      -- Opaque surfaces keep each variant readable regardless of terminal colors.
      palette = {
        dawn = {
          muted = '#716b80',
          subtle = '#655f75',
        },
        moon = {
          base = '#2a273f',
          muted = '#aaa6bf',
          subtle = '#c0bcd1',
          surface = '#393552',
          overlay = '#44415a',
          highlight_low = '#393552',
          highlight_med = '#56526e',
          highlight_high = '#6e6a86',
        },
      },
      highlight_groups = {
        LspInlayHint = { bg = 'base', fg = 'muted', italic = true },
        NotificationInfo = { bg = 'none', fg = 'text' },
        NotificationWarning = { bg = 'none', fg = 'subtle' },
        NotificationError = { bg = 'none', fg = 'love' },
      },
    }
    vim.cmd 'colorscheme rose-pine'
  end,
}
