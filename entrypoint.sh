#!/bin/bash

# Install bash

# Unpack args into an array
args=("$@")
# command_args=()

# Get command, path and include-preview
command="$1"
path="$2"
if [ "$3" = "--include-preview=true" ]; then
    include_preview="--include-preview"
elif [ "$3" = "--include-preview=false" ]; then
    include_preview=""
else
    echo "Invalid value for --include-preview (true/false)"
    exit 1
fi

# Remove first 3 args
unset 'args[0]'
unset 'args[0]'
unset 'args[0]'

# Print args
echo "Command: $command"
echo "Path: $path"
echo "Include preview: $include_preview"

for i in "${args[@]}"; do
    echo "Arg: $i"
done

# # If command is "scan"
# if [ "$command" = "scan" ]; then
    
#     for i in "${args[@]}"; do
#         # if arg starts with "--include-preview"


#     done

# elif [ "$command" = "connect" ]; then
#     for i in "${args[@]}"; do
#         echo "Arg: $i"
#     done

# else 
#     echo "Command not found"
#     exit 1
# fi
