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
- Go

## Tests

Tests are kept in the root `tests/` directory and run with:

```bash
make test
```

## Setup

```bash
cp .env.example .env
docker compose up -d db
make migrate-up
docker compose up --build
```

## Quiz data API

Categories group questions but do not determine their answer format. Each question
references its own quiz type. Use `GET /quiz-types` to retrieve quiz type IDs,
then include `quizTypeId` when creating a question:

```json
{
  "question": "What is the capital of Germany?",
  "categoryId": 1,
  "quizTypeId": 2,
  "answers": [
    { "answer": "Berlin", "isCorrect": true },
    { "answer": "Hamburg", "isCorrect": false }
  ]
}
```

Create categories with only a name, for example `{"name":"Geography"}`.

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