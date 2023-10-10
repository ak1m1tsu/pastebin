include .env.example
export

LOCAL_BIN=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: deps
deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag-v1:
	swag init -g internal/controller/http/v1/router.go

.PHONY: run
run: swag-v1
	go mod tidy && go mod download && \
	go run -tags migrate ./cmd/app
