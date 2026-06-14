# User Age API

A RESTful API built with **Go**, **GoFiber**, **PostgreSQL + SQLC**, **Uber Zap**, and
**go-playground/validator** that manages users with a `name` and `dob` (date of birth),
dynamically calculating each user's `age` using Go's `time` package.

## Project Structure

```
/cmd/server/main.go        - application entrypoint
/config/                    - configuration loading (env vars)
/db/migrations/              - SQL migration files
/db/queries/                  - SQL queries used by sqlc
/db/sqlc/                     - sqlc-generated DB access code
/internal/
├── handler/                  - HTTP handlers (Fiber)
├── repository/                - data access layer (wraps sqlc)
├── service/                    - business logic (age calculation, etc.)
├── routes/                      - route registration
├── middleware/                   - request ID + request logging middleware
├── models/                        - request/response DTOs + age calc + tests
└── logger/                         - Zap logger setup
```

## Prerequisites

- Go 1.22+
- PostgreSQL 13+ (or Docker)
- [sqlc](https://docs.sqlc.dev/) (only needed if you change queries/schema)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (optional, for manual migrations)
- Docker & Docker Compose (optional, recommended)

## Quick Start (Docker — recommended)

```bash
docker compose up --build
```

This will:
1. Start a PostgreSQL container
2. Run database migrations automatically
3. Build and start the API on `http://localhost:3000`

## Manual Setup

1. **Clone the repo and install dependencies**

   ```bash
   git clone <your-repo-url>
   cd user-age-api
   go mod download
   ```

2. **Start PostgreSQL** (locally or via Docker)

   ```bash
   docker run --name user-age-db -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=user_age_db \
     -p 5432:5432 -d postgres:16-alpine
   ```

3. **Run migrations**

   ```bash
   migrate -path db/migrations \
     -database "postgres://postgres:postgres@localhost:5432/user_age_db?sslmode=disable" up
   ```

4. **Configure environment variables**

   Copy `.env.example` to `.env` and adjust as needed, or export the variables directly:

   ```bash
   cp .env.example .env
   export $(cat .env | xargs)
   ```

5. **Run the server**

   ```bash
   go run ./cmd/server
   ```

   The API will be available at `http://localhost:3000`.

## Running Tests

```bash
go test ./... -v
```

This includes a dedicated unit test for the `CalculateAge` function in `internal/models/user_test.go`.

## API Endpoints

| Method | Endpoint      | Description                          |
|--------|---------------|---------------------------------------|
| POST   | `/users`      | Create a new user                     |
| GET    | `/users`      | List users (paginated)                |
| GET    | `/users/:id`  | Get a single user (with calculated age)|
| PUT    | `/users/:id`  | Update a user                         |
| DELETE | `/users/:id`  | Delete a user (returns 204)           |

### Create User

```http
POST /users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

**Response (201):**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID

```http
GET /users/1
```

**Response (200):**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

### List Users (paginated)

```http
GET /users?page=1&page_size=10
```

**Response (200):**
```json
{
  "data": [
    { "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 1,
  "total_pages": 1
}
```

### Update User

```http
PUT /users/1
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

**Response (200):**
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User

```http
DELETE /users/1
```

**Response:** `204 No Content`

## Features

- ✅ Clean layered architecture (handler → service → repository → sqlc)
- ✅ Age calculated dynamically with Go's `time` package (no stored age)
- ✅ Input validation via `go-playground/validator`
- ✅ Structured logging via Uber Zap (including per-request logs)
- ✅ `X-Request-ID` middleware for request tracing
- ✅ Request duration logging middleware
- ✅ Pagination on `GET /users`
- ✅ Unit tests for age calculation
- ✅ Docker & Docker Compose support
- ✅ Proper HTTP status codes (`201`, `200`, `204`, `400`, `404`, `500`)

## Environment Variables

| Variable     | Default       | Description              |
|--------------|---------------|---------------------------|
| `DB_HOST`    | `localhost`   | PostgreSQL host            |
| `DB_PORT`    | `5432`        | PostgreSQL port            |
| `DB_USER`    | `postgres`    | PostgreSQL user            |
| `DB_PASSWORD`| `postgres`    | PostgreSQL password        |
| `DB_NAME`    | `user_age_db` | PostgreSQL database name   |
| `APP_PORT`   | `3000`        | Port the API listens on    |

## Regenerating SQLC Code

If you modify `db/queries/users.sql` or `db/migrations/`, regenerate the sqlc code with:

```bash
sqlc generate
```
