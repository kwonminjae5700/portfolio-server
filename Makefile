.PHONY: help build run test clean docker-up docker-down migrate

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build files
	rm -rf bin/
	rm -f coverage.out

docker-up: ## Start docker containers
	docker-compose up -d

docker-down: ## Stop docker containers
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

docker-rebuild: ## Rebuild and restart docker containers
	docker-compose down
	docker-compose up -d --build

migrate: ## Run database migrations
	go run ./cmd/migrate/main.go

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
