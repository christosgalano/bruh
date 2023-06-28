#!/bin/bash

# Function to extract flag
extract_flag() {
  if [[ "$1" == *=false ]]; then
    echo ""
  elif [[ "$1" == *=true ]]; then
    echo "${1%=true}"
  else
    return 1
  fi
}

# Get command, path and include-preview
command="$1"
path="$2"

include_preview=$(extract_flag "$3")
return_code=$?
if [[ $return_code -eq 1 ]]; then
  echo "Error: Invalid argument for --include-preview (true/false))"
  exit 1
fi

echo "Command: $command"
echo "Path: $path"
echo "Include preview: $include_preview"

# Get the appropriate arguments for the command
if [[ "$command" == "scan" ]]; then
    outdated=$(extract_flag "$4")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --outdated (true/false))"
      exit 1
    fi
    output="$5"

    echo "Outdated: $outdated"
    echo "Output: $output"

elif [[ "$command" == "update" ]]; then
    in_place=$(extract_flag "$6")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --in-place (true/false))"
      exit 1
    fi
    silent=$(extract_flag "$7")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --silent (true/false))"
      exit 1
    fi

    echo "In place: $in_place"
    echo "Silent: $silent"

else 
    echo "Error: Command not found (scan/update)"
    exit 1
fi
