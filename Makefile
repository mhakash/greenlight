.PHONY: run

run:
	@export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable'; \
	go run ./cmd/api
