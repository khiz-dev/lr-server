
POSTGRES_USER := "postgres"
POSTGRES_PASSWORD := "postgres"
POSTGRES_DB := "chirpy"
POSTGRES_PORT := "5432"
POSTGRES_URL := "postgres://" + POSTGRES_USER + ":" + POSTGRES_PASSWORD + "@localhost:" + POSTGRES_PORT + "/" + POSTGRES_DB + "?sslmode=disable"

# Start PostgreSQL 15
run_database:
    #!/usr/bin/env bash
    set -e
    if docker container inspect postgres15 >/dev/null 2>&1; then
        docker start postgres15
    else
        docker run -d \
            --name postgres15 \
            -e POSTGRES_USER={{POSTGRES_USER}} \
            -e POSTGRES_PASSWORD={{POSTGRES_PASSWORD}} \
            -e POSTGRES_DB={{POSTGRES_DB}} \
            -p 127.0.0.1:{{POSTGRES_PORT}}:5432 \
            -v postgres15_data:/var/lib/postgresql/data \
            --restart unless-stopped \
            postgres:15
    fi

# Stop PostgreSQL
stop_database:
    docker stop postgres15

# Connect to PostgreSQL
connect_psql:
    psql "{{POSTGRES_URL}}"

# Apply migrations
goose_up:
    goose -dir sql/schema postgres "{{POSTGRES_URL}}" up

# Rollback migration
goose_down:
    goose -dir sql/schema postgres "{{POSTGRES_URL}}" down

# Show migration status
goose_status:
    goose -dir sql/schema postgres "{{POSTGRES_URL}}" status

# Run the server with the local database configuration
run:
    DB_URL="{{POSTGRES_URL}}" ENV=dev go run .

show_env:
    echo "DB_URL={{POSTGRES_URL}}"
    echo "ENV=dev"