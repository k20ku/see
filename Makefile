.DEFAULT_GOAL := help

DOKERTAG := latest

.PHONY: build
build: ## build doker image to deploy

	docker build -t k20ku/see:${DOKERTAG} \
		--target deploy ./

.PHONY: build-local
build-local: ## build docker image to local development

	docker compose build --no-cache app-dev

.PHONY: up
up: ## Do docker compose up

	docker compose up -d app-dev

.PHONY: down
down: ## Do docker compose down

	docker compose down

.PHONY: logs
logs: ## Tail docker compose logs

	docker compose logs -f

.PHONY: ps
ps: ## Check container status

	docker compose ps

.PHONY: tests
.tests: ## Execute tests

	go test -race -shuffle=on ./...


# Code Quality
.PHONY: fmt
fmt: ## Format codes

	golangci-lint fmt

.PHONY: lint
lint: ## Lint code

	golangci-lint run

.PHONY: check
check: fmt lint tests ## Check code quality

.PHONY: help
help: ## Show options
	
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
