# Go Packing Service

API-only Go/Gin service that calculates optimal pack allocations and persists configurable pack sizes in PostgreSQL with optimistic concurrency.

## Tech Stack

- Go
- Gin HTTP framework
- PostgreSQL (`database/sql` + `lib/pq`, no ORM)
- Viper (JSON profile config loading)
- `slog` JSON logging
- Docker + Docker Compose
- Go standard testing library

## Configuration (Viper)

Configuration is loaded by environment profile:

- `APP_ENV=dev` -> `configs/dev.json`
- `APP_ENV=prod` -> `configs/prod.json`

If `APP_ENV` is not set, `dev` is used.

### Config files

- `configs/dev.json`
- `configs/prod.json`

### Environment overrides

These environment variables override file values:

- `APP_ENV`
- `PORT` (maps to `server.port`)
- `DATABASE_URL` (maps to `database.url`)

## VS Code Launch Configs

VS Code launch profiles are in `.vscode/launch.json`:

- `Go Packing API (dev)`
- `Go Packing API (prod)`

Each profile sets `APP_ENV` accordingly.

## Run with Docker Compose

```bash
docker compose up --build
```

The compose app service runs with `APP_ENV=prod`, so it reads `configs/prod.json` by default.
