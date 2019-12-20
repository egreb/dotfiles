set runtimepath+=~/.vim_runtime

" install plugins
call plug#begin()
Plug 'junegunn/seoul256.vim'
Plug 'majutsushi/tagbar'
Plug 'junegunn/fzf', { 'dir': '~/.fzf', 'do': './install --all' }
Plug 'kien/ctrlp.vim'
Plug 'scrooloose/nerdcommenter'
call plug#end()

" Change mapleader
let mapleader=","

" color scheme
set t_Co=256
let g:seoul256_background = 234
colo seoul256

" general stuff
set tabstop=4
set number
set autoindent					" Copy indent from last line when starting new line.
set autowrite
filetype plugin on
set backspace=indent,eol,start " Make sure backspace works as intended.

" keybinding
" map next item in quicklist to ctrl+n
map <C-n> :cnext<CR>
" map previous item in quicklist to ctrl+m
map <C-m> :cprevious<CR>
" close quicklist with , then a
nnoremap <leader>a :cclose<CR>

" nerdcomment: github.com/scrooloose/nerdcommenter
let g:NERDSpaceDelims = 1
let g:NERDCommentEmptyLines = 1
let g:NERDTrimTrailingWhitespace = 1
let g:NERDTrimTrailingWhitespace = 1
