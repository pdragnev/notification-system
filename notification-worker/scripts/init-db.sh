#!/bin/bash
# Navigate to the script's directory
cd "$(dirname "$0")"

# Attempt to load environment variables from .env file located in the parent directory
if [ -f "../.env" ]; then
    echo "Loading environment variables from .env file..."
    export $(grep -v '^#' ../.env | xargs)
fi

export PGPASSWORD=$DB_PASS

echo "Creating database '$DB_NAME'..."

psql -U $DB_USER -h $DB_HOST -d postgres -c "CREATE DATABASE $DB_NAME;"

echo "Initializing database '$DB_NAME'..."

psql -U $DB_USER -h $DB_HOST -d $DB_NAME -a -f ./init-users.sql

echo "Database '$DB_NAME' initialized successfully."