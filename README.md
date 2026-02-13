# Go Packing Service

Run the API and PostgreSQL locally with Docker Compose.

## Start locally

```bash
docker compose -p packaging up -d
```

What this starts:

- `go-packing-postgres` on `localhost:5400`
- `go-packing-api` on `localhost:8080`
- `go-packing-ui` on `localhost:3000`

Swagger UI:

- `http://localhost:8080/swagger/index.html`

Web UI:

- `http://localhost:3000`
