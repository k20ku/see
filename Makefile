.DEFAULT_GOAL := help

DOKERTAG := latest
GOBIN := $(CURDIR)/tools/bin

# Build

.PHONY: build
build: ## build docker image to local development
	docker compose build --no-cache --profile dev

.PHONY: build-deploy
build-deploy: ## build doker image to deploy
	docker build -t k20ku/see:${DOKERTAG} \
		--target deploy ./

.PHONY: up
up: ## Do docker compose up for development
	docker compose --profile dev up

.PHONY: down
down: ## Do docker compose down
	docker compose --profile dev down

.PHONY: logs
logs: ## Tail docker compose logs
	docker compose logs -f

.PHONY: ps
ps: ## Check container status
	docker compose ps -a

# Tool Install

.PHONY: tools
tools: ## Install tools
	@mkdir -p $(GOBIN)
	GOBIN=$(GOBIN) go install tool
	GOBIN=$(GOBIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(GOBIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "install golangci-lint"
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOBIN) v2.12.2

.PHONY: migrate
migrate: ## Migration
	$(GOBIN)/migrate -database "postgres://see:seedbpass@localhost:54321/see?sslmode=disable" -path db/migrations up

.PHONY: immigrate
immigrate: ## Immigration
	$(GOBIN)/migrate -database "postgres://see:seedbpass@localhost:54321/see?sslmode=disable" -path db/migrations down

.PHONY: sqlgen
sqlgen: ## Generate GO code from SQL
	$(GOBIN)/sqlc generate

# Code Quality

.PHONY: tests
tests: ## Execute tests
	go test -race -shuffle=on ./...

.PHONY: fmt
fmt: ## Format codes
	GOBIN=$(GOBIN) golangci-lint fmt

.PHONY: lint
lint: ## Lint code
	GOBIN=$(GOBIN) golangci-lint run

.PHONY: check
check: fmt lint tests ## Check code quality

# Help
.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
