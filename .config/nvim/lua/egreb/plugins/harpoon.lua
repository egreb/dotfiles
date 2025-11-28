return {
	"ThePrimeagen/harpoon",
	branch = "harpoon2",
	dependencies = { "nvim-lua/plenary.nvim" },
	config = function()
		local harpoon = require('harpoon')
		harpoon:setup()

		vim.keymap.set("n", "<leader>fa", function() harpoon:list():add() end)
		vim.keymap.set("n", "<C-e>", function() harpoon.ui:toggle_quick_menu(harpoon:list()) end)

		-- Toggle previous & next buffers stored within Harpoon list
		vim.keymap.set("n", "<leader>fp", function() harpoon:list():prev() end)
		vim.keymap.set("n", "<leader>fn", function() harpoon:list():next() end)
	end
}
