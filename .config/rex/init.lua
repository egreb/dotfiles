-- Rex configuration. Apply changes with `rex config reload`.

package.path = os.getenv("HOME") .. "/code/sbx/rex/?.lua;" .. package.path
package.loaded["sbx"] = nil -- load the latest sbx.lua on every rex config reload

require("sbx").setup({
    project_roots = { "~/code" }, -- folders cmd+f searches for git repos
})

-- Work notes: each action opens Neovim in a centered float and runs an `:Obsidian` subcommand.
-- Quit Neovim to close the float.
local function notes_action(name, title, subcommand)
    rex.action({
        name = name,
        title = title,
        keywords = "obsidian vault nvim",
        run = function()
            local _, err = rex.session.new_layer({
                bounds = { 0.05, 0.05, 0.9, 0.9 },
                focus = true,
                layout = rex.layout.block({
                    flavor = "com.superlogical.terminal.shell",
                    options = {
                        -- Starting in the vault makes obsidian.nvim pick the work workspace
                        cwd = os.getenv("HOME") .. "/vaults/work",
                        command = { "nvim", "-c", "Obsidian " .. subcommand },
                    },
                }),
            })
            if err ~= nil then
                error(err, 0)
            end
        end,
    })
end

notes_action("notes", "Open Work Notes", "quick_switch") -- note picker
notes_action("notes_new", "New Work Note", "new") -- prompts for a title

rex.bind("shift+cmd+o", "notes")
rex.bind("alt+cmd+n", "notes_new")
