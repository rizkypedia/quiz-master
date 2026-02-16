DB_URL=postgres://app:app@localhost:5432/quiz_master?sslmode=disable

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1

migrate-force:
	docker compose run --rm migrate force $(version)

migrate-version:
	docker compose run --rm migrate version
