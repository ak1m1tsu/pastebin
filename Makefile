include .envrc
export

LOCAL_BIN=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## deps: install migrate and swag utils to ./bin
.PHONY: deps
deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## swag: generate swagger documentation for http handlers
.PHONY: swag
swag-v1:
	swag fmt && swag init -g internal/controller/http/v1/router.go

## run: start the cmd/app application
.PHONY: run/app
run/app: swag-v1
	go run ./cmd/app

## migrations/new name=$1: create a new migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path="./migrations" -database ${POSTGRES_DSN} up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path="./migrations" -database ${POSTGRES_DSN} down

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker/up: build and up docker services
.PHONY: docker/up
docker/up: confirm
	docker-compose up --build -d && docker-compose logs -f

## docker/down: down and remove docker services
.PHONY: docker/down
docker/down: confirm
	docker-compose down --remove-orphans

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code 
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependecies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependecies...'
	go mod vendor
	@echo 'Formatting code...'
	gofumpt -w ./
	@echo 'Linting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/app: build the cmd/app application
.PHONY: build/app
build/app:
	@echo 'Building cmd/app...'
	go build -ldflags '-s -w' -o ./bin/app ./cmd/app
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o=./bin/linux_amd64/app ./cmd/app
