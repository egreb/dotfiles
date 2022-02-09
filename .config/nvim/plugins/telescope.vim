Plug 'nvim-lua/popup.nvim'
Plug 'nvim-lua/plenary.nvim'
Plug 'nvim-telescope/telescope.nvim'

" if !exists('g:loaded_telescope') | finish | endif

nnoremap <leader>f <cmd>Telescope find_files<cr>
nnoremap <leader>p <cmd>Telescope git_files<cr>
nnoremap <leader>r <cmd>Telescope live_grep<cr>
nnoremap <leader>b <cmd>Telescope buffers<cr>
nnoremap <silent>; <cmd>Telescope help_tags<cr>

