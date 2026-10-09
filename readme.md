# Quiz Master API

Quiz Master is a backend API for building quiz applications. It provides endpoints
for managing quiz categories, questions, answer choices, and users, so a frontend
or another client can create and deliver quizzes without implementing its own
quiz-data service.

Categories group related questions. Each question has its own quiz type, allowing
single-choice and multiple-choice questions to coexist in the same category.
Questions contain answer choices, with correct answers identified by the API.
User registration and login are supported with JWT-based authentication.

## Technology stack

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

The backend provides the API and database used by the
[Quiz Master frontend](https://github.com/rizkypedia/quiz-master-fe). The
frontend is a separate project and cannot load categories, run quizzes, or
calculate results without this backend running.

```bash
cp .env.example .env
# Set FRONTEND_URL=http://localhost:5174 in .env for the local Vite frontend.
docker compose up -d db
make migrate-up
docker compose up --build
```

The API listens on `http://localhost:8080`. Start this backend and its database
before starting the frontend. In a separate terminal, follow the install and
run instructions in the
[frontend README](https://github.com/rizkypedia/quiz-master-fe#readme); the
frontend's Vite server uses `http://localhost:5174` by default. `FRONTEND_URL`
must match the frontend's address: the backend uses it both to generate
shareable quiz links and to allow browser requests through CORS.

## Frontend integration

The frontend is the participant-facing UI; the backend remains the source of
truth for category availability, quiz questions, correct answers, secure links,
and scores. The frontend calls this API directly and uses these endpoints:

- `GET /categories` to show active categories.
- `POST /quiz-link` with `{"categoryId":1}` to create a secure share link.
- `GET /quiz?token=...` to retrieve the questions for that category. A valid
  secure token is required.
- `POST /quiz/result?token=...` to submit answers and receive the score.

Secure links expire after 12 hours. Deactivating a category prevents new links
and invalidates access through its existing links. The quiz API returns answer
choices without correctness information; scoring is performed by the backend.

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

Create categories with only a name, for example `{"name":"Geography"}`. New
categories are active by default. `GET /categories` returns active categories
only; authenticate and send `PUT /category/:id/active` with
`{"active":false}` or `{"active":true}` to change availability. Inactive
categories cannot be selected for new quiz links, and links to them stop working.

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