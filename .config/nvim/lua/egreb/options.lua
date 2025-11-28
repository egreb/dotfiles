-- set leader to <space>
vim.g.mapleader        = ' '
vim.g.maplocalleader   = ' '

vim.opt.number         = true
vim.opt.relativenumber = true
vim.opt.tabstop        = 4
vim.opt.shiftwidth     = 4
vim.opt.swapfile       = false
vim.opt.wrap           = false
vim.opt.signcolumn     = 'yes'
vim.opt.winborder      = 'rounded'
vim.opt.smartindent    = true
vim.opt.termguicolors  = true
vim.opt.ignorecase     = true
vim.opt.cursorcolumn   = false
vim.opt.cursorline     = true
vim.opt.undofile       = true
vim.g.have_nerd_font   = true

-- Decrease update time
vim.opt.updatetime     = 250

-- Decrease mapped sequence wait time
-- Displays which-key popup sooner
vim.opt.timeoutlen     = 300

-- Don't show the mode, since it's already in the status line
vim.opt.showmode       = false

-- Sync clipboard between OS and Neovim.
--  Remove this option if you want your OS clipboard to remain independent.
--  See `:help 'clipboard'`
vim.opt.clipboard      = 'unnamedplus'

-- Minimal number of screen lines to keep above and below the cursor.
vim.o.scrolloff        = 10

-- if performing an operation that would fail due to unsaved changes in the buffer (like `:q`),
-- instead raise a dialog asking if you wish to save the current file(s)
-- See `:help 'confirm'`
vim.o.confirm          = true
vim.opt.smoothscroll   = true
vim.opt.foldtext       = ""
vim.wo.foldmethod      = 'expr'
vim.wo.foldexpr        = 'v:lua.vim.treesitter.foldexpr()'
vim.opt.foldlevel      = 99
vim.opt.foldlevelstart = 1
