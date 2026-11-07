-- In-buffer markdown rendering (headings, callouts, checkboxes, tables, wikilinks).
-- Owns the markdown UI; obsidian.nvim's built-in UI is disabled in obsidian.lua.
return {
  'MeanderingProgrammer/render-markdown.nvim',
  ft = { 'markdown' },
  dependencies = {
    'nvim-treesitter/nvim-treesitter', -- markdown + markdown_inline already installed
    'echasnovski/mini.nvim', -- provides mini.icons (enabled in mini.nvim.lua)
  },
  ---@module 'render-markdown'
  ---@type render.md.UserConfig
  opts = {
    completions = { lsp = { enabled = true } }, -- callout/checkbox completion via blink's lsp source
  },
}
