migrate:
	cd db/migrations; \
	GOOSE_DBSTRING=postgresql://postgres:postgres@localhost:5432/postgres go tool goose postgres up

generate:
	go tool sqlc generate

generate-mocks:
	mkdir -p mocks
	go generate ./internal/port/...

lint:
	go tool golangci-lint run

clean-mocks:
	rm -rf mocks

.PHONY: migrate generate generate-mocks lint clean-mocks