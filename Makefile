migrate:
	cd db/migrations; \
	GOOSE_DBSTRING=postgresql://postgres:postgres@localhost:5432/postgres go tool goose postgres up

generate:
	go tool sqlc generate
