#!/bin/bash

# Script: sqlc-generate-delete-query.sh
# Usage: ./sqlc-generate-delete-query.sh <table_name> <primary_key> <output_file>

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Error: Missing required parameters"
    echo "Usage: $0 <table_name> <primary_key> <output_file>"
    exit 1
fi

TABLE_NAME=$1
PRIMARY_KEY=$2
OUTPUT_FILE=$3

# Create the output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Convert TABLE_NAME to PascalCase and remove trailing 's' if exists
FORMATTED_NAME=$(echo "$TABLE_NAME" | awk -F'_' '{for (i=1; i<=NF; i++) printf "%s", toupper(substr($i,1,1)) tolower(substr($i,2))}' | sed 's/s$//')

# Generate the DELETE query
QUERY_CONTENT="-- name: Delete${FORMATTED_NAME} :exec
DELETE FROM ${TABLE_NAME}
WHERE ${PRIMARY_KEY} = $1;"

# Check if file exists
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "$QUERY_CONTENT" > "$OUTPUT_FILE"
    echo "✅ Created new SQL file: $OUTPUT_FILE"
    exit 0
fi

# Check if query with same name already exists
if grep -qE "^-- name: Delete${FORMATTED_NAME} :exec" "$OUTPUT_FILE"; then
    echo "❌ Error: Query 'Delete${FORMATTED_NAME}' already exists in $OUTPUT_FILE"
    exit 1
fi

# Append new query to existing file
echo "$QUERY_CONTENT" >> "$OUTPUT_FILE"
echo "✅ Appended new query to: $OUTPUT_FILE" 