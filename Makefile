.PHONY: docker-up docker-down docker-down-v migrate-up migrate-down

docker-up:
	docker compose up -d postgres

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1
