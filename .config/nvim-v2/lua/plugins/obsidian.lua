-- Use the stable release API with Snacks and blink's existing LSP source.
vim.pack.add({
	{ src = "https://github.com/obsidian-nvim/obsidian.nvim", version = "v3.16.7" },
})

local opts = {
    workspaces = {
      { name = 'work', path = '~/vaults/work' },
    },
    -- Match Obsidian's built-in daily notes: plain files named YYYY-MM-DD in the
    -- vault root. Omitting `folder` keeps them in the root; workdays_only = false
    -- uses literal calendar days (no weekend skipping) like Obsidian.
    daily_notes = {
      date_format = '%Y-%m-%d',
      default_tags = {},
      workdays_only = false,
      template = nil,
    },
    -- Use the picker that's already enabled in snacks.lua.
    picker = {
      name = 'snacks.picker',
    },
    -- render-markdown.nvim owns the in-buffer UI, so disable obsidian's built-in renderer.
    ui = {
      enable = false,
    },
    -- Completion is provided via an in-process LSP and surfaces through blink's `lsp` source.
    completion = {
      min_chars = 2,
      create_new = true,
    },
    -- Keep obsidian.nvim's managed frontmatter (id/aliases/tags) on regular notes,
    -- but disable it for daily notes so they stay plain like Obsidian's. `fname` is
    -- the vault-relative path; daily notes are `YYYY-MM-DD.md` in the vault root.
    frontmatter = {
      enabled = function(fname)
        return not (fname ~= nil and fname:match '^%d%d%d%d%-%d%d%-%d%d%.md$' ~= nil)
      end,
    },
}

require('obsidian').setup(opts)

local keys = {
    { '<leader>mt', '<cmd>Obsidian today<cr>', desc = 'Obsidian: today (daily note)' },
    { '<leader>my', '<cmd>Obsidian yesterday<cr>', desc = 'Obsidian: yesterday note' },
    { '<leader>mn', '<cmd>Obsidian new<cr>', desc = 'Obsidian: new note' },
    { '<leader>ms', '<cmd>Obsidian search<cr>', desc = 'Obsidian: search notes' },
    { '<leader>mq', '<cmd>Obsidian quick_switch<cr>', desc = 'Obsidian: quick switch' },
    { '<leader>mb', '<cmd>Obsidian backlinks<cr>', desc = 'Obsidian: backlinks' },
    { '<leader>ml', '<cmd>Obsidian links<cr>', desc = 'Obsidian: links in note' },
    { '<leader>mg', '<cmd>Obsidian tags<cr>', desc = 'Obsidian: tags (search / add via <C-x>/<C-l>)' },
}

for _, key in ipairs(keys) do
	vim.keymap.set("n", key[1], key[2], { desc = key.desc })
end
