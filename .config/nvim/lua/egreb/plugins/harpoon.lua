return {
  'ThePrimeagen/harpoon',
  branch = 'harpoon2',
  dependencies = { 'nvim-lua/plenary.nvim' },
  keys = {
    {
      '<leader>fa',
      function()
        require('harpoon'):list():add()
      end,
      desc = 'Harpoon add file',
    },
    {
      '<leader>fe',
      function()
        local harpoon = require 'harpoon'
        harpoon.ui:toggle_quick_menu(harpoon:list())
      end,
      desc = 'Harpoon quick menu',
    },
    -- Toggle previous & next buffers stored within Harpoon list
    {
      '<leader>fp',
      function()
        require('harpoon'):list():prev()
      end,
      desc = 'Harpoon prev',
    },
    {
      '<leader>fn',
      function()
        require('harpoon'):list():next()
      end,
      desc = 'Harpoon next',
    },
  },
  config = function()
    require('harpoon'):setup()
  end,
}
