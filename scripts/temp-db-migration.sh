#!/bin/bash

# Check if the VERSION variable is not provided
if [ -z "$1" ]; then
  echo "Error: VERSION is not provided."
  echo "Usage: ./create_migration.sh <VERSION>"
  exit 1
fi

VERSION="$1"
touch ./db/migrations/${VERSION}_.sql;
# The name of the migration file
FILE="./db/migrations/${VERSION}_.sql"

# Create the file and insert the SQL content
cat <<EOL > "$FILE"
-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd
-- Migration missing from auto-generated[automatic-from-make-command]
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
EOL

# Successfully notify
echo "Migration file created: $FILE"
