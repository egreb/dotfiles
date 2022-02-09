export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# powerlevel10k
if [ -f ~/.dotfiles/shell/powerlevel10k ]; then
	source ~/.dotfiles/shell/powerlevel10k
else
	print "powerlevel10k config not found"
fi

# HISTORY
if [ -f ~/.dotfiles/shell/history ]; then
	source ~/.dotfiles/shell/history
else
	print "404: history config file not found"
fi

# aliases
if [ -f ~/.dotfiles/shell/alias ]; then
	source ~/.dotfiles/shell/alias
else 
	print "404: alias file not found"
fi

# zsh / oh-my-zsh config
if [ -f ~/.dotfiles/shell/zsh ]; then
	source ~/.dotfiles/shell/zsh
else 
	print "404: zsh / oh-my-zsh config file not found"
fi

# paths
if [ -f ~/.dotfiles/shell/paths ]; then
	source ~/.dotfiles/shell/paths
else
	print "404: paths file not found"
fi

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
