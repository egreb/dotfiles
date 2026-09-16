-- [[ Basic Autocommands ]]
--  See `:help lua-guide-autocommands`

-- Highlight when yanking (copying) text
--  Try it with `yap` in normal mode
--  See `:help vim.hl.on_yank()`
vim.api.nvim_create_autocmd('TextYankPost', {
  desc = 'Highlight when yanking (copying) text',
  group = vim.api.nvim_create_augroup('egreb-highlight-yank', { clear = true }),
  callback = function()
    vim.hl.on_yank()
  end,
})

-- restore cursor pos on file open
vim.api.nvim_create_autocmd('BufReadPost', {
  pattern = '*',
  callback = function()
    local line = vim.fn.line '\'"'
    if line > 1 and line <= vim.fn.line '$' then
      vim.cmd 'normal! g\'"'
    end
  end,
})

-- disable automatic comment on newline
vim.api.nvim_create_autocmd('FileType', {
  pattern = '*',
  callback = function()
    vim.opt_local.formatoptions:remove { 'c', 'r', 'o' }
  end,
})

-- follow the macOS light/dark appearance
--  Ghostty swaps its own theme via `theme = light:...,dark:...`, but it has no way
--  to tell nvim, so ask macOS directly. Checked on focus rather than on a timer:
--  flipping appearance always means leaving and re-entering the terminal.
local function sync_appearance()
  vim.system({ 'defaults', 'read', '-g', 'AppleInterfaceStyle' }, { text = true }, function(obj)
    local want = (obj.stdout or ''):match 'Dark' and 'dark' or 'light'
    vim.schedule(function()
      if vim.o.background ~= want then
        vim.o.background = want
        vim.cmd.colorscheme 'rose-pine'
      end
    end)
  end)
end

vim.api.nvim_create_autocmd({ 'VimEnter', 'FocusGained' }, {
  desc = 'Match colorscheme to the macOS light/dark appearance',
  group = vim.api.nvim_create_augroup('egreb-appearance', { clear = true }),
  callback = sync_appearance,
})
