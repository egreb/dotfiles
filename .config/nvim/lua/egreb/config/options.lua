-- Leader key (needed by all `<leader>` keybindings)
vim.g.mapleader = " " -- space is the leader key

-- Core UI2 (experimental): replaces the builtin message + cmdline presentation layer
require("vim._core.ui2").enable({
	enable = true,
	msg = {
		-- where each kind of message is shown: "cmd", "msg" (floating) or "pager"
		targets = {
			[""] = "msg",
			empty = "cmd",
			bufwrite = "msg",
			confirm = "cmd",
			emsg = "pager",
			echo = "msg",
			echomsg = "msg",
			echoerr = "pager",
			completion = "cmd",
			list_cmd = "pager",
			lua_error = "pager",
			lua_print = "msg",
			progress = "pager",
			rpc_error = "pager",
			quickfix = "msg",
			search_cmd = "cmd",
			search_count = "cmd",
			shell_cmd = "pager",
			shell_err = "pager",
			shell_out = "pager",
			shell_ret = "msg",
			undo = "msg",
			verbose = "pager",
			wildlist = "cmd",
			wmsg = "msg",
			typed_cmd = "cmd",
		},
		dialog = { height = 0.5 }, -- max height of modal prompts
		msg = { height = 0.3 }, -- max height of the message window
		pager = { height = 0.5 }, -- max height of the pager window
	},
})
vim.opt.messagesopt:append("timeout:5000") -- ms a message stays visible in the ui2 message window
vim.opt.messagesopt:append("maxheight:50") -- max height of the expanded cmdline, as a % of 'lines'

-- Colorscheme (builtin)
vim.o.termguicolors = true -- enable 24-bit colors
vim.cmd.colorscheme("catppuccin") -- load the builtin catppuccin theme

-- UI
vim.o.number = true -- show the absolute line number on the current line
vim.o.relativenumber = true -- show other line numbers relative to the cursor
vim.o.cursorline = true -- highlight the line the cursor is on
vim.o.wrap = false -- don't wrap long lines, scroll horizontally instead
vim.o.scrolloff = 10 -- keep 10 lines visible above/below the cursor
vim.o.signcolumn = "yes" -- always show the sign column so text doesn't shift
vim.o.cmdheight = 0 -- hide the command line when it's not in use
vim.o.showmatch = true -- briefly highlight the matching bracket when one is inserted
vim.o.matchtime = 2 -- how long showmatch highlights, in tenths of a second
vim.o.pumheight = 10 -- max number of items in the popup menu
vim.o.pumblend = 10 -- popup menu transparency (0-100)
vim.o.winborder = "rounded" -- rounded borders on floating windows
vim.o.splitbelow = true -- horizontal splits open below the current window
vim.o.splitright = true -- vertical splits open to the right of the current window

-- Cursor (experimental): block in normal/visual/command, thin bar in insert, underline in replace
vim.opt.guicursor =
	"n-v-c:block,i-ci-ve:ver25,r-cr:hor20,o:hor50,a:blinkwait700-blinkoff400-blinkon250-Cursor/lCursor,sm:block-blinkwait175-blinkoff150-blinkon175"

-- Tabs: 4 spaces
vim.o.tabstop = 4 -- a tab character is displayed as 4 columns wide
vim.o.softtabstop = 4 -- pressing <Tab>/<BS> moves 4 columns
vim.o.shiftwidth = 4 -- indent with 4 columns for >>, << and auto-indent
vim.o.expandtab = true -- insert spaces instead of tab characters

-- Search
vim.o.ignorecase = true -- ignore case when searching
vim.o.smartcase = true -- ...unless the pattern contains an uppercase letter
vim.o.wildignore = "node_modules" -- skip node_modules in file completion

-- Files
vim.o.clipboard = "unnamedplus" -- yank/paste use the system clipboard
vim.o.swapfile = false -- don't create swap files
vim.o.undofile = true -- persist undo history across sessions
local undodir = vim.fn.expand("~/.cache/nvim/undodir")
vim.fn.mkdir(undodir, "p") -- make sure the undo directory exists
vim.o.undodir = undodir -- where undo history is stored

-- Timing
vim.o.updatetime = 300 -- ms of idle time before CursorHold fires
vim.o.timeoutlen = 500 -- ms to wait for a mapped key sequence to complete
vim.o.ttimeoutlen = 0 -- ms to wait for a key code sequence (no delay leaving insert mode)

-- Input
vim.o.mouse = "" -- disable the mouse

-- Netrw
vim.g.netrw_keepdir = 0 -- keep the browsing directory in sync with the current directory
vim.g.netrw_localcopydircmd = "cp -r" -- copy directories recursively

-- Diagnostics
vim.diagnostic.config({
	virtual_text = {
		prefix = "●", -- marker shown before the message
		spacing = 2, -- spaces between the code and the message
		source = "if_many", -- show source (e.g. gopls) only when several are attached
		current_line = false, -- virtual_lines below takes over on the current line
	},
	-- full, wrapped message rendered below the line the cursor is on
	virtual_lines = { current_line = true },
	underline = true, -- underline the code that has a diagnostic
	severity_sort = true, -- show the most severe diagnostics first
	float = { border = "rounded", source = true }, -- rounded float showing the source
})
