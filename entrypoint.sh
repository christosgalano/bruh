#!/bin/bash

# Install bash

# Unpack args into an array
args=("$@")

command="$1"
unset 'array[0]'

# Print args
echo "Command: $command"
for i in "${args[@]}"; do
    echo "Arg: $i"
done

# If command is "scan"
if [ "$command" = "scan" ]; then
    echo "Scanning..."
elif [ "$command" = "connect" ]; then
    echo "Connecting..."
else 
    echo "Command not found"
    exit 1
fi
