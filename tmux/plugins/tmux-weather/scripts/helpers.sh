#!/usr/bin/env bash

PATH="/usr/local/bin:$PATH:/usr/sbin"

use_current_location() {
  echo -n node -pe 'JSON.parse(process.argv[2]).city' "$(curl -s ipinfo.io)"
}

get_tmux_option() {
  local option_name="$1"
  local default_value="$2"
  local option_value=$(tmux show-option -gqv $option_name)

  if [ -z "$option_value" ]; then
    echo -n $default_value
  else
    echo -n $option_value
  fi
}

set_tmux_option() {
  local option_name="$1"
  local option_value="$2"
  $(tmux set-option -gq $option_name "$option_value")
}

