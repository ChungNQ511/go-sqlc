#!/bin/bash

# Script: sqlc-generate-insert-query.sh
# Usage: ./sqlc-generate-insert-query.sh <table_name> <columns> <output_file>

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Error: Missing required parameters"
    echo "Usage: $0 <table_name> <columns> <output_file>"
    exit 1
fi

TABLE_NAME=$1
COLUMNS=$2
OUTPUT_FILE=$3

# Create the output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Convert TABLE_NAME to PascalCase and remove trailing 's' if exists
FORMATTED_NAME=$(echo "$TABLE_NAME" | awk -F'_' '{for (i=1; i<=NF; i++) printf "%s", toupper(substr($i,1,1)) tolower(substr($i,2))}' | sed 's/s$//')

# Format columns with proper indentation
FORMATTED_COLUMNS=$(echo "$COLUMNS" | awk -F',' '{for(i=1;i<=NF;i++) printf "    %s%s\n", $i, (i==NF?"":",")}')
PLACEHOLDERS=$(echo "$COLUMNS" | awk -F',' '{for(i=1;i<=NF;i++) printf "    $%d%s\n", i, (i==NF?"":",")}')

# Generate the INSERT query
QUERY_CONTENT="-- name: Create${FORMATTED_NAME} :one
INSERT INTO ${TABLE_NAME} (
${FORMATTED_COLUMNS}
) VALUES (
${PLACEHOLDERS}
) RETURNING *;"

# Check if file exists
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "$QUERY_CONTENT" > "$OUTPUT_FILE"
    echo "✅ Created new SQL file: $OUTPUT_FILE"
    exit 0
fi

# Check if query with same name already exists
if grep -qE "^-- name: Create${FORMATTED_NAME} :one" "$OUTPUT_FILE"; then
    echo "❌ Error: Query 'Create${FORMATTED_NAME}' already exists in $OUTPUT_FILE"
    exit 1
fi

# Append new query to existing file
echo "$QUERY_CONTENT" >> "$OUTPUT_FILE"
echo "✅ Appended new query to: $OUTPUT_FILE" 