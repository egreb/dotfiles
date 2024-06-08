#!/bin/bash

if [ -n "$1" ]; then
    selected_repo="$1"
else
    # Capture the selected repository name from fzf into a variable
    selected_repo=$(gh repo ls cybr-ai | while IFS= read -r w; do
    echo "$w" | awk '{split($1, a, "/"); print a[2]}'
    done | fzf)
fi

# Check if a repository was selected
if [ -n "$selected_repo" ]; then
    # Original URL with placeholder
    url="https://grafana.stage.apps.cybr.ai/explore?schemaVersion=1&panes=%7B%22lz3%22%3A%7B%22datasource%22%3A%22P8E80F9AEF21F6940%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22expr%22%3A%22%7Bapp%3D%5C%22harambe-executor%5C%22%7D%22%2C%22queryType%22%3A%22range%22%2C%22datasource%22%3A%7B%22type%22%3A%22loki%22%2C%22uid%22%3A%22P8E80F9AEF21F6940%22%7D%2C%22editorMode%22%3A%22code%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-15m%22%2C%22to%22%3A%22now%22%7D%7D%7D&orgId=1"

    # Replace 'harambe-executor' with the selected repository name
    updated_url=$(echo "$url" | sed "s/harambe-executor/$selected_repo/")

    # Print the updated URL
    open $updated_url
else
    echo "No repository selected."
fi

