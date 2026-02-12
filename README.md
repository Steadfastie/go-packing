# Go Packing Service

Run the API and PostgreSQL locally with Docker Compose.

## Start locally

```bash
docker compose -p packaging up -d --build
```

What this starts:

- `go-packing-postgres` on `localhost:5400`
- `go-packing-db-init` to create database/table from `docker/postgres/init.sql`
- `go-packing-api` on `localhost:8080`

Swagger UI:

- `http://localhost:8080/swagger/index.html`
