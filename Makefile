.PHONY: fmt vet lint test run build docker-up docker-down docker-logs \
        db-up db-down db-reset db-logs \
        sqlc sqlc-watch \
        migrate-up migrate-down migrate-force migrate-create migrate-reset \
        dev-init dev-reset check test-log

include .env
export

DATETIME = $(shell powershell -Command "Get-Date -format 'yyyy-MM-dd_HH-mm-ss'")
MIGRATIONS_PATH=internal/adapter/postgres/migrations

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test -v ./... -count=1

test-log:
	@if not exist logs mkdir logs
	go test -v ./... -count=1 > logs/test_$(DATETIME).log

run:
	go run cmd/app/main.go

build: test swagger
	go build -o bin/app cmd/app/main.go

# --- DOCKER ---
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# --- POSTGRES ---
db-up:
	docker compose up -d postgres

db-down:
	docker compose stop postgres

db-reset:
	docker compose down -v
	docker compose up -d postgres
	ping 127.0.0.1 -n 4 > NUL

db-logs:
	docker compose logs -f postgres

db-seed:
	psql "$(DATABASE_URL)" -f ./db_seed.sql
	@echo "Seeding completed"

# --- SQLC ---
sqlc:
	sqlc generate

sqlc-watch:
	find . -name "*.sql" | entr -r make sqlc

# --- MIGRATIONS ---
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-force:
	@read -p "version? " v; \
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $$v

migrate-create:
	@read -p "migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name

migrate-reset:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

# --- DEV FLOWS ---
dev-init: db-up migrate-up sqlc
	@echo "Dev environment ready"

dev-reset: db-reset migrate-up sqlc
	@echo "Dev environment reset"

pre-commit:
	pre-commit run --all-files

swagger:
	swag init -g cmd/app/main.go