# Setup the application
setup:
	@chmod +x setup.sh
	@./setup.sh

# Run the application locally
run: setup
	@go run ./cmd/api

# Build the binary locally
build:
	@echo "Building binary..."
	@go build -o main ./cmd/api

# Clean built files
clean:
	@echo "Cleaning up..."
	@rm -f main

# Run tests
test:
	@echo "Running tests..."
	@go test -count=1 ./...

# Live Reload with Air
watch: setup
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Docker: Build containers
docker-build: setup
	@docker compose build

# Docker: Build and start all containers in background
docker-run: setup
	@docker compose up --build -d

# Docker: Stop all containers
docker-down:
	@docker compose down
