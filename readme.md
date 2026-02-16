# Go Gin API Starter

This is a template for Apps with Go Gin

- Go 1.25
- Gin
- PostgreSQL
- JWT authentication
- golang-migrate
- Docker & Docker Compose
- Hot reload via Air

## Requirements

- Docker
- Docker Compose
- Make

## Setup

```bash
cp .env.example .env
docker compose up -d db
make migrate-up
docker compose up --build
```

## Migration Restart from scratch (Clean DB)


```
docker compose down -v

docker compose up -d

make migrate-up
```


## Migration Fix dirty DB

Only when

- You want to keep your database volume AND
- You want to fix a dirty migration AND
- You do not want to lose data

````
make migrate-force version=1
make migrate-down
make migrate-up