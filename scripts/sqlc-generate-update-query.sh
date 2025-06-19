#!/bin/bash

# Script: sqlc-generate-update-query.sh
# Usage: ./sqlc-generate-update-query.sh <table_name> <columns> <primary_key> <output_file>

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ]; then
    echo "Error: Missing required parameters"
    echo "Usage: $0 <table_name> <columns> <primary_key> <output_file>"
    exit 1
fi

TABLE_NAME=$1
COLUMNS=$2
PRIMARY_KEY=$3
OUTPUT_FILE=$4

# Create the output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Convert TABLE_NAME to PascalCase and remove trailing 's' if exists
FORMATTED_NAME=$(echo "$TABLE_NAME" | awk -F'_' '{for (i=1; i<=NF; i++) printf "%s", toupper(substr($i,1,1)) tolower(substr($i,2))}' | sed 's/s$//')

# Format columns with proper indentation and placeholders
FORMATTED_COLUMNS=$(echo "$COLUMNS" | awk -F',' '{for(i=1;i<=NF;i++) printf "    %s = $%d%s\n", $i, i, (i==NF?"":",")}')

# Generate the UPDATE query
QUERY_CONTENT="-- name: Update${FORMATTED_NAME} :one
UPDATE ${TABLE_NAME}
SET ${FORMATTED_COLUMNS}
WHERE ${PRIMARY_KEY} = $(( $(echo "$COLUMNS" | tr -cd ',' | wc -c) + 1 ))
RETURNING *;"

# Check if file exists
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "$QUERY_CONTENT" > "$OUTPUT_FILE"
    echo "✅ Created new SQL file: $OUTPUT_FILE"
    exit 0
fi

# Check if query with same name already exists
if grep -qE "^-- name: Update${FORMATTED_NAME} :one" "$OUTPUT_FILE"; then
    echo "❌ Error: Query 'Update${FORMATTED_NAME}' already exists in $OUTPUT_FILE"
    exit 1
fi

# Append new query to existing file
echo "$QUERY_CONTENT" >> "$OUTPUT_FILE"
echo "✅ Appended new query to: $OUTPUT_FILE" 