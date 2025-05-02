#!/bin/bash

set -e

JSON_FILE="./test_data/${1}"

# Check for jq dependency
if ! command -v jq &> /dev/null; then
  echo "Error: 'jq' is not installed." >&2
  exit 1
fi

jq -c '.[]' "$JSON_FILE" | while read -r row; do
  rowkey=$(echo "$row" | jq -r '.rowkey')
  family=$(echo "$row" | jq -r '.family')

  # Start command
  cmd=(go run . write -f "$family" -k "$rowkey")

  # Loop over qualifiers
  qualifiers=$(echo "$row" | jq -c '.qualifiers')
  for key in $(echo "$qualifiers" | jq -r 'keys[]'); do
    value=$(echo "$qualifiers" | jq -r --arg k "$key" '.[$k]')
    cmd+=(-q "$key" -v "$value")
  done

  # Print and run the command
  echo "Running: ${cmd[*]}"
  "${cmd[@]}"
done
