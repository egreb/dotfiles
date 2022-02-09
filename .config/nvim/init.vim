autocmd!
scriptencoding utf-8
if !1 | finish | endif

set nocompatible
set number
syntax enable
set fileencodings=utf-8,sjis,latin
set title
set autoindent
set background=dark
set nobackup
set hlsearch
set showcmd
set cmdheight=1
set laststatus=2
set scrolloff=10
set expandtab
set shell=fish
set backupskip=/tmp/*,/private/tmp/*
set splitright
set relativenumber
set ignorecase
set smarttab
set shiftwidth=2
set tabstop=2
set ai "autoindent
set si "smartindent
set nowrap "No wrap lines
set backspace=start,eol,indent
set path+=**
set wildignore+=*/node_modules/*
set formatoptions+=r
set cursorline

filetype plugin indent on
autocmd InsertLeave * set nopaste

if has('nvim')
	set inccommand=split
endif

"--------------------------------------------------------------------------
" Key maps
"--------------------------------------------------------------------------

let mapleader = "\<space>"

nmap <leader>ve :edit ~/.config/nvim/init.vim<cr>
nmap <leader>vr :source ~/.config/nvim/init.vim<cr>

nmap <leader>k :nohlsearch<CR>
nmap <leader>Q :bufdo bdelete<cr>

" Allow gf to open non-existent files
map gf :edit <cfile><cr>

let data_dir = has('nvim') ? stdpath('data') . '/site' : '~/.vim'
call plug#begin(data_dir . '/plugins')
source ~/.config/nvim/plugins/telescope.vim
source ~/.config/nvim/plugins/commenary.vim
source ~/.config/nvim/plugins/seoul256.vim
call plug#end()

colo seoul256
