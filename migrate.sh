#!/bin/bash
source .env

construct_migration_dsn() {
    local db_host="$DB_HOST"
    local db_port="$DB_PORT"
    local db_name="$POSTGRES_DB"
    local db_user="$POSTGRES_USER"
    local db_password="$POSTGRES_PASSWORD"
    local ssl_mode="disable"

    echo "host=$db_host port=$db_port dbname=$db_name user=$db_user password=$db_password sslmode=$ssl_mode"
}

MIGRATION_DSN=$(construct_migration_dsn)

echo "Running migrations..."

sleep 2 && goose -dir "${MIGRATIONS_DIR}" postgres "${MIGRATION_DSN}" up -v
