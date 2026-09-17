vim.pack.add({
  "https://github.com/zhisme/copy_with_context.nvim",
})

require("copy_with_context").setup({
  mappings = {
    relative = "<leader>cy",
    absolute = "<leader>cY",
    remote = "<leader>cr",
  },
  formats = {
    default = "# {filepath}:{line}",
    remote = "# {remote_url}",
  },
  trim_lines = false,
})
