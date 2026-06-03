# Simple Makefile for a Go project

setup:
	@chmod +x setup.sh
	@./setup.sh

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Docker — app only (uses external postgres, set host in config.yaml)
docker-run:
	@docker compose up --build -d

# Docker — app + postgres (set host: postgres in config.yaml)
docker-run-all:
	@docker compose --profile db up --build -d

docker-down:
	@docker compose --profile db down

docker-logs:
	@docker compose logs -f app

.PHONY: all build run test clean watch docker-run docker-run-all docker-down docker-logs itest
