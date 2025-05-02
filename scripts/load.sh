#!/bin/bash

set -e

JSON_FILE="./test_data/${1}"
MAX_PARALLEL=50 # Adjust this based on your system capabilities

# Check for dependencies
if ! command -v jq &> /dev/null; then
  echo "Error: 'jq' is not installed." >&2
  exit 1
fi

if ! command -v parallel &> /dev/null; then
  echo "Error: 'parallel' is not installed." >&2
  echo "Install with: brew install parallel (macOS)" >&2
  exit 1
fi

# Create a temp file to store commands
CMDS_FILE=$(mktemp)
trap 'rm -f $CMDS_FILE' EXIT

# Generate all commands first
jq -c '.[]' "$JSON_FILE" | while read -r row; do
  rowkey=$(echo "$row" | jq -r '.rowkey')
  family=$(echo "$row" | jq -r '.family')

  cmd="go run . write -f \"$family\" -k \"$rowkey\""

  # Loop over qualifiers
  qualifiers=$(echo "$row" | jq -c '.qualifiers')
  for key in $(echo "$qualifiers" | jq -r 'keys[]'); do
    value=$(echo "$qualifiers" | jq -r --arg k "$key" '.[$k]')
    cmd+=" -q \"$key\" -v \"$value\""
  done

  # Add command to file
  echo "$cmd" >> "$CMDS_FILE"
done

# Execute commands in parallel
cat "$CMDS_FILE" | parallel -j $MAX_PARALLEL --verbose "bash -c {}"
