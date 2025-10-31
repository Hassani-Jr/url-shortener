.PHONY: run test build docker-up docker-down clean

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Build binary
build:
	go build -o bin/api cmd/api/main.go

# Start Docker services
docker-up:
	docker compose up -d

# Stop Docker services
docker-down:
	docker compose down

# Clean build artifacts
clean:
	rm -rf bin/ tmp/ coverage.out

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Database commands
db-create:
	docker exec -it urlshortener-db createdb -U postgres urlshortener

db-drop:
	docker exec -it urlshortener-db dropdb -U postgres urlshortener

db-shell:
	docker exec -it urlshortener-db psql -U postgres -d urlshortener