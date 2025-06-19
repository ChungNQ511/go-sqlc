#!/bin/bash

MIGRATIONS_DIR="./db/migrations"

# The name of the migration to check
MIGRATION_NAME="$1"

# Check if no parameter is provided
if [ -z "$MIGRATION_NAME" ]; then
  echo "Please provide the migration name. Example: ./create-migration.sh create_table_abc"
  exit 1
fi

# Check if the migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "The directory $MIGRATIONS_DIR does not exist. Please check again."
  exit 1
fi

# Check if the migration file already exists
if ls "$MIGRATIONS_DIR"/*"$MIGRATION_NAME".sql 1> /dev/null 2>&1; then
  echo "❌ Migration with the name '$MIGRATION_NAME' already exists in the directory $MIGRATIONS_DIR."
  exit 1
else
  echo "✅ Creating a new migration with the name '$MIGRATION_NAME'..."
  goose -dir $MIGRATIONS_DIR create "$MIGRATION_NAME" sql
fi