.DEFAULT_GOAL := help

DOKERTAG := latest

.PHONY: build
build: ## build doker image to deploy

	docker build -t k20ku/see:${DOKERTAG} \
		--target deploy ./

.PHONY: build-local
build-local: ## build docker image to local development

	docker compose build --no-cache

.PHONY: up
up: ## Do docker compose

	docker compose up -d

.PHONY: logs
logs: ## Tail docker compose logs

	docker compose logs -f

.PHONY: ps
ps: ## Check container status

	docker compose ps

.PHONY: tests
.tests: ## Execute tests

	go test -race -shuffle=on ./...

.PHONY: help
help: ## Show options
	
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
