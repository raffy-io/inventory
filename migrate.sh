#!/bin/bash

set -e

# Automatically load .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default to "up" if no argument is provided
COMMAND="${1:-up}"
TARGET_VERSION="${2:-0}"

# Falls back to a safe default if DB_URL isn't found in .env
DB_URL="${DB_URL:-postgresql://postgres:your_password@localhost:5432/your_db_name}"
MIGRATION_DIR="database/migration"

if [ "$COMMAND" = "down-to-zero" ] || [ "$COMMAND" = "reset" ]; then
    echo "Resetting database (migrating down to version 0)..."
    goose -dir "$MIGRATION_DIR" postgres "$DB_URL" down-to "$TARGET_VERSION"
else
    echo "Running goose migrations ($COMMAND)..."
    goose -dir "$MIGRATION_DIR" postgres "$DB_URL" "$COMMAND"
fi