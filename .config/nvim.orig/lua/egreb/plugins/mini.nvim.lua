return { -- Collection of various small independent plugins/modules
  'echasnovski/mini.nvim',
  name = 'mini',
  version = false,
  config = function()
    -- Better Around/Inside textobjects
    --
    -- Examples:
    --  - va)  - [V]isually select [A]round [)]paren
    --  - yinq - [Y]ank [I]nside [N]ext [']quote
    --  - ci'  - [C]hange [I]nside [']quote
    require('mini.ai').setup {}

    -- Add/delete/replace surroundings (brackets, quotes, etc.)
    -- - saiw) - [S]urround [A]dd [I]nner [W]ord [)]Paren
    -- - sd'   - [S]urround [D]elete [']quotes
    -- - sr)'  - [S]urround [R]eplace [)] [']
    require('mini.surround').setup {}
    -- require('mini.indentscope').setup {}
    require('mini.move').setup {}
    -- require('mini.pick').setup {}
    -- require('mini.pairs').setup {}
    require('mini.jump').setup {}

    -- Git diff signs in the gutter (add/change/delete hunks).
    -- style = 'sign' forces gutter signs; the default would tint line numbers since `number` is on.
    -- Also adds hunk mappings: gh (apply), gH (reset), [h/]h (prev/next hunk), gh (textobject).
    require('mini.diff').setup {
      view = {
        style = 'sign',
        signs = { add = '▎', change = '▎', delete = '▁' },
      },
    }

    -- Icon provider used by render-markdown.nvim (headings, callouts, etc.)
    require('mini.icons').setup {}
  end,
}
