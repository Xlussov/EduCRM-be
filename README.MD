# EduCRM Backend

Backend system for an educational center to manage branches, teachers, students, lesson schedules, attendance, and subscription plans. Built with Go and PostgreSQL.

## Prerequisites

- Go 1.25.7
- Docker and Docker Compose
- Make

*Note: Required CLI dependencies for local development (`golang-migrate`, `sqlc`, `swag`) are installed automatically via the initialization script. Run `make setup` (which executes `sh scripts/sh/init.sh`) to install them.*

## Configuration

The application requires both environment variables and a YAML configuration file to run correctly.

### 1. Environment Variables (.env)
Create a `.env` file in the project root directory based on the provided template:
```bash
cp example.env .env
```
Ensure the following core variables are defined in your `.env` file:
- `CONFIG_PATH`: Path to the configuration file (e.g., `config/local.yaml`).
- `DATABASE_URL`: PostgreSQL connection string (used by Make, migrations, and the application).
- JWT configurations (`JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`).
- PostgreSQL parameters (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB_NAME`, `POSTGRES_PORT`).

### 2. YAML Configuration (config.yaml)
Application settings not managed by environment secrets must reside in the YAML file specified by the `CONFIG_PATH` environment variable. This includes the application environment state and server network parameters.

Example configuration structure:
```yaml
env: "local" # Available options: local, dev, prod
http_server:
  address: ":8080"
  read_timeout: 5s
  write_timeout: 10s
```
Ensure the file exists at the specified path and is properly configured before starting the server.

## Running the Application

The server starts by default at `http://localhost:8080`.

### Option A: Fully Containerized (Docker Compose)
To run both the application and the PostgreSQL database in containers:
```bash
make docker-up
```

### Option B: Local Development
To run the database in Docker and the application locally:
1. Start the PostgreSQL container:
   ```bash
   make db-up
   ```
2. Apply migrations and generate SQLC models:
   ```bash
   make dev-init
   ```
3. Start the application:
   ```bash
   make run
   ```

## Database Initialization & Seeding

The system requires an initial `SUPERADMIN` account to access the API. Test data can also be injected via SQL scripts.

1. **Create Superadmin**:
   ```bash
   make db-create-superadmin
   ```
2. **Seed Initial Data** (Branches, Subjects, Students, Teachers):
   ```bash
   make db-seed
   ```

## API Documentation

The API uses Swagger for documentation. To generate or update the documentation:
```bash
make swagger
```
Once the application is running, the Swagger UI is accessible at:
`http://localhost:8080/docs`

## Makefile Command Reference

A comprehensive list of all available `make` commands to manage the project lifecycle.

### Build & Run
| Command | Description |
|---|---|
| `make run` | Starts the application locally (`go run cmd/app/main.go`). |
| `make build` | Runs tests, generates swagger, and builds the binary to `bin/app`. |
| `make setup` | Executes the initialization script located at `scripts/sh/init.sh`. |

### Docker Operations
| Command | Description |
|---|---|
| `make docker-up` | Builds and starts all containers in detached mode. |
| `make docker-down` | Stops and removes all containers. |
| `make docker-logs` | Tails logs from all running containers. |

### PostgreSQL Operations
| Command | Description |
|---|---|
| `make db-up` | Starts only the PostgreSQL container. |
| `make db-down` | Stops the PostgreSQL container. |
| `make db-reset` | Removes database volumes, restarts the container, and verifies connection. |
| `make db-logs` | Tails logs specifically from the PostgreSQL container. |
| `make db-seed` | Executes `db_seed.sql` against the database. |
| `make db-create-superadmin`| Executes `create_superadmin.sql` to initialize the root user. |

### SQLC (Code Generation)
| Command | Description |
|---|---|
| `make sqlc` | Generates Go code from SQL queries using `sqlc`. |
| `make sqlc-watch` | Watches `.sql` files and automatically regenerates code on changes. |

### Migrations
| Command | Description |
|---|---|
| `make migrate-up` | Applies all pending migrations to the database. |
| `make migrate-down` | Rolls back the last applied migration (down 1). |
| `make migrate-force` | Forces a specific migration version (prompts for version number). |
| `make migrate-create`| Creates a new set of up/down migration files (prompts for name). |
| `make migrate-reset` | Rolls back all migrations and applies them again from scratch. |

### Development Workflow
| Command | Description |
|---|---|
| `make dev-init` | Starts DB container, applies migrations, and runs sqlc generation. |
| `make dev-reset` | Resets the DB container completely, applies migrations, and runs sqlc. |

### Code Quality & Testing
| Command | Description |
|---|---|
| `make fmt` | Formats the codebase using `go fmt`. |
| `make vet` | Runs `go vet` to identify potential issues. |
| `make lint` | Executes `golangci-lint run`. |
| `make test` | Runs all unit tests verbosely without caching (`-count=1`). |
| `make pre-commit` | Manually triggers `pre-commit` hooks on all files. |
| `make swagger` | Updates Swagger documentation using `swag init`. |

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

© 2026 Danil Xlussov
