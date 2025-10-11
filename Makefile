include .env
MIGRATION_NAME ?= m

run api:
	@go run cmd/url-shortener/main.go

migration_sql:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./storage/migrations create $(MIGRATION_NAME) sql
migration_go:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./storage/migrations create $(MIGRATION_NAME) go
migrate_up:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./storage/migrations up
migrate_down:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./storage/migrations down
migrate_status:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./storage/migrations status
