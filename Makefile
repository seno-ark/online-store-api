#!make
include .env

dev:
	go run cmd/api/main.go

swagger:
	swag init --parseDependency --parseInternal --dir "cmd/api,internal/api" --output internal/api/docs

migrate-file:
	migrate create -ext sql --dir pkg/database/migration -seq $(name)

migrate-up:
	migrate -path pkg/database/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migrate-down:
	migrate -path pkg/database/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

.PHONY: dev swagger migrate-file migrate-up migrate-down