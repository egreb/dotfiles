set runtimepath+=~/.vim_runtime

call plug#begin()
Plug 'junegunn/seoul256.vim'
"Plug 'majutsushi/tagbar'
"Plug 'fatih/vim-go', { 'tag': '*', 'do': ':GoInstallBinaries' }
"Plug 'nsf/gocode', { 'rtp': 'vim', 'do': '~/.vim/plugged/gocode/vim/symlink.sh' }
Plug 'junegunn/fzf', { 'dir': '~/.fzf', 'do': './install --all' }
Plug 'kien/ctrlp.vim'
call plug#end()

"let g:deoplete#enable_at_startup = 1
set t_Co=256
let g:seoul256_background = 234
colo seoul256

set backspace=indent,eol,start

let g:deoplete#enable_at_startup = 1 " ignored by vim because deoplete not installed there
" set autowrite
" map next item in quicklist to ctrl+n
" map <C-n> :cnext<CR>
" map previous item in quicklist to ctrl+m
" map <C-m> :cprevious<CR>
" close quicklist with , then a
" nnoremap <leader>a :cclose<CR>
" build file with , b
" autocmd FileType go nmap <leader>b  <Plug>(go-build)
" run file with , r
" autocmd FileType go nmap <leader>r  <Plug>(go-run)
" use goimports
" let g:go_fmt_command = "goimports""

filetype plugin on
