#!/bin/bash

# Directory to watch
WATCH_DIR="/"

# Commands to run on change
COMMAND1="../go/bin/templ generate"
COMMAND2="go run main.go"

# Using fswatch to monitor for changes
fswatch -o "$WATCH_DIR" | while read change; do
    echo "Change detected, running commands..."
    
    # Run the commands
    $COMMAND1
    $COMMAND2
done