.PHONY: run

run:
	@export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable'; \
	go run ./cmd/api

migrate-up:
	@export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable'; \
	migrate -path ./migrations -database $$GREENLIGHT_DB_DSN up

migrate-create:
	@echo "Enter the name of the migration:"
	@read -r migration_name; \
	migrate create -seq -ext=.sql -dir=./migrations $$migration_name