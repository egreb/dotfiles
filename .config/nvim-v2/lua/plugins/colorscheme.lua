vim.pack.add({
  "https://github.com/maxmx03/solarized.nvim",
})

vim.opt.termguicolors = true
vim.opt.background = "light"

require("solarized").setup({
  transparent = {
    enabled = false,
  },

  styles = {
    enabled = true,
    comments = {},
    keywords = {},
    functions = {},
    variables = {},
  },

  on_highlights = function(colors)
    return {
      Comment = {
        fg = colors.base0,
      },

      CursorLine = {
        bg = colors.base2,
      },

      Visual = {
        bg = colors.base2,
      },

      Search = {
        fg = colors.base02,
        bg = colors.yellow,
      },

      IncSearch = {
        fg = colors.base3,
        bg = colors.orange,
      },

      DiagnosticUnderlineError = {
        sp = colors.red,
        underline = true,
      },

      DiagnosticUnderlineWarn = {
        sp = colors.yellow,
        underline = true,
      },

      DiagnosticUnderlineInfo = {
        sp = colors.blue,
        underline = true,
      },

      DiagnosticUnderlineHint = {
        sp = colors.cyan,
        underline = true,
      },
    }
  end,
})

vim.cmd.colorscheme("solarized")
