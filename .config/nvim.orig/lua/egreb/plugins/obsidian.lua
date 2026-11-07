return {
  'obsidian-nvim/obsidian.nvim',
  version = '*',
  lazy = true,
  ft = 'markdown',
  dependencies = {
    'nvim-lua/plenary.nvim',
  },
  keys = {
    { '<leader>mt', '<cmd>Obsidian today<cr>', desc = 'Obsidian: today (daily note)' },
    { '<leader>my', '<cmd>Obsidian yesterday<cr>', desc = 'Obsidian: yesterday note' },
    { '<leader>mn', '<cmd>Obsidian new<cr>', desc = 'Obsidian: new note' },
    { '<leader>ms', '<cmd>Obsidian search<cr>', desc = 'Obsidian: search notes' },
    { '<leader>mq', '<cmd>Obsidian quick_switch<cr>', desc = 'Obsidian: quick switch' },
    { '<leader>mb', '<cmd>Obsidian backlinks<cr>', desc = 'Obsidian: backlinks' },
    { '<leader>ml', '<cmd>Obsidian links<cr>', desc = 'Obsidian: links in note' },
    { '<leader>mg', '<cmd>Obsidian tags<cr>', desc = 'Obsidian: tags (search / add via <C-x>/<C-l>)' },
    { '<leader>mw', '<cmd>Obsidian workspace<cr>', desc = 'Obsidian: switch workspace' },
  },
  ---@module 'obsidian'
  ---@type obsidian.config
  opts = {
    workspaces = {
      { name = 'personal', path = '~/vaults/personal' },
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
    -- NOTE: the enum value is 'snacks.pick' (not 'snacks.picker') in the installed
    -- v3.x release; it's also kept as a backward-compat alias on newer versions.
    picker = {
      name = 'snacks.pick',
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
  },
  config = function(_, opts)
    -- Remember the last-selected workspace across sessions. obsidian.nvim only
    -- picks a workspace from the cwd (or the first spec) with no memory between
    -- sessions, so we persist it to a small state file and restore it on load.
    local state_file = vim.fs.joinpath(vim.fn.stdpath 'state', 'obsidian-last-workspace')

    -- Read the saved name *before* setup(): setup()'s initial Workspace.set() fires
    -- the callback below and would otherwise overwrite the saved value first.
    local saved
    if vim.fn.filereadable(state_file) == 1 then
      saved = vim.trim(vim.fn.readfile(state_file)[1] or '')
      if saved == '' then saved = nil end
    end

    -- Persist the workspace whenever it changes (startup or `:Obsidian workspace`).
    opts.callbacks = opts.callbacks or {}
    local user_cb = opts.callbacks.post_set_workspace
    opts.callbacks.post_set_workspace = function(workspace)
      pcall(vim.fn.writefile, { workspace.name }, state_file)
      if user_cb then user_cb(workspace) end
    end

    require('obsidian').setup(opts)

    -- Restore the saved workspace, but only when the cwd isn't already inside a
    -- vault (there, cwd-based detection is intentional and should win).
    if saved then
      local ws = require 'obsidian.workspace'
      local in_vault = ws.find(vim.uv.cwd(), Obsidian.workspaces)
      if not in_vault and Obsidian.workspace.name ~= saved then
        pcall(ws.set, saved)
      end
    end
  end,
}
