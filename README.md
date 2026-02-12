# Go Packing Service

API-only Go/Gin service that calculates optimal packs for an order amount while enforcing these rules:

1. Only whole packs can be shipped.
2. Minimize total shipped quantity first.
3. Within that minimum quantity, minimize number of packs.

The service also provides PostgreSQL-backed pack size configuration with optimistic concurrency.

## Tech Stack

- Go
- Gin HTTP framework
- PostgreSQL (raw SQL via `database/sql` + `pgx` stdlib)
- `log/slog` JSON logging
- Docker + Docker Compose
- Go standard testing package

## Project Structure

```text
cmd/api/main.go                             # Application entrypoint and dependency wiring
internal/presentation/http/                 # HTTP handlers, router, middleware
internal/service/                           # Use-case orchestration
internal/domain/                            # Domain models, errors, solver
internal/infrastructure/postgres/           # PostgreSQL access layer
pkg/logx/                                   # Shared JSON logger bootstrap
pkg/httpx/                                  # Shared API error helpers
docker/postgres/init.sql                    # DB schema init script
Dockerfile                                  # App container image
docker-compose.yml                          # Local app + postgres stack
```

## API

Base URL: `http://localhost:8080`

### Health Check

`GET /healthz`

Response:

```json
{"status":"ok"}
```

### Calculate Packs

`POST /api/v1/calculate`

Request:

```json
{
  "amount": 500000
}
```

Success response (`200`): returns only packs array.

```json
[
  { "size": 53, "count": 9429 },
  { "size": 31, "count": 7 },
  { "size": 23, "count": 2 }
]
```

Errors:

- `400` invalid request body or invalid amount
- `409` pack sizes not configured
- `500` internal error

### Get Pack Sizes

`GET /api/v1/pack-sizes`

Response when not configured yet:

```json
{
  "version": 0,
  "pack_sizes": []
}
```

Response when configured:

```json
{
  "version": 3,
  "pack_sizes": [23, 31, 53]
}
```

### Replace Pack Sizes

`PUT /api/v1/pack-sizes`

Request:

```json
{
  "pack_sizes": [23, 31, 53]
}
```

Success response (`200`):

```json
{
  "version": 4,
  "pack_sizes": [23, 31, 53]
}
```

Errors:

- `400` invalid `pack_sizes` (empty, non-positive, duplicates)
- `409` optimistic concurrency conflict
- `500` internal error

## Optimistic Concurrency

The API does not accept a version in `PUT` requests. Instead:

1. Service loads current configuration (`version = N`).
2. Domain `Replace()` validates and increments to `N+1`.
3. Persistence executes CAS update with condition `db.version = N`.
4. If no row matches, repository returns conflict and API responds `409`.

## Running with Docker Compose

```bash
docker compose up --build
```

Services:

- App: `localhost:8080`
- PostgreSQL: `localhost:5432`

Database schema is created by Postgres init script at startup:

- `docker/postgres/init.sql`

No migration framework is used.

## Usage Example

1. Try calculate before configuration (expect `409`):

```bash
curl -i -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"amount":251}'
```

2. Configure pack sizes:

```bash
curl -i -X PUT http://localhost:8080/api/v1/pack-sizes \
  -H 'Content-Type: application/json' \
  -d '{"pack_sizes":[250,500,1000,2000,5000]}'
```

3. Calculate again:

```bash
curl -i -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"amount":12001}'
```

Expected result includes:

- `2 x 5000`
- `1 x 2000`
- `1 x 250`

## Edge Case Verification

Configure edge case pack sizes:

```bash
curl -X PUT http://localhost:8080/api/v1/pack-sizes \
  -H 'Content-Type: application/json' \
  -d '{"pack_sizes":[23,31,53]}'
```

Calculate for `500000`:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"amount":500000}'
```

Expected output:

- `{23: 2, 31: 7, 53: 9429}`

(JSON form):

```json
[
  { "size": 53, "count": 9429 },
  { "size": 31, "count": 7 },
  { "size": 23, "count": 2 }
]
```

## Unit Tests

Run tests (inside containerized Go toolchain):

```bash
docker run --rm -v "$PWD:/src" -w /src golang:1.22 sh -lc "go test ./..."
```

## Logging

The service uses `slog` JSON logging. Each request is logged with method, path, status, and duration.

## Git History

Changes were implemented as incremental commits to show development progression instead of one large commit.
