# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Neovim configuration using Lua, organized under the `egreb` namespace. Part of a larger dotfiles repo.

## Architecture

**Bootstrap**: `init.lua` requires five modules in order:
1. `egreb.options` — vim options (leader=space, 4-space indent, relative numbers)
2. `egreb.lazy` — lazy.nvim plugin manager setup (auto-bootstraps from GitHub)
3. `egreb.keymaps` — global keybindings
4. `egreb.autocommands` — autocmds (yank highlight, cursor restore, no auto-comment)
5. `egreb.lsp` — LSP keybindings and diagnostic config

**Plugin organization**: Each plugin gets its own file in `lua/egreb/plugins/`. LSP-related plugins are in `lua/egreb/plugins/lsp/`. lazy.nvim loads specs from both directories.

**After directory**: `after/lsp/` contains per-server LSP configs (ts_ls, emmet_ls). `after/queries/go/` has TreeSitter injection queries for SQL in Go strings.

## Key Design Decisions

- **blink.cmp** for completion (not nvim-cmp)
- **snacks.nvim** as unified picker/UI framework (not telescope)
- **conform.nvim** for formatting with format-on-save enabled (prettier for web, stylua for lua, goimports+gofmt for go)
- **mason.nvim** manages LSP servers: ts_ls, html, tailwindcss, lua_ls, emmet_ls, jsonls, bashls, vue_ls
- **sidekick.nvim** for AI assistant integration (Claude/Copilot via tmux backend)
- `install.missing = false` in lazy.nvim — plugins won't auto-install

## Lua Formatting

Uses stylua (configured in `.stylua.toml`): 160 column width, 2-space indent, single quotes, no call parentheses.
