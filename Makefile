.PHONY: build run test clean tidy

APP_NAME=fluffy-mouton-api
MAIN_FILE=cmd/api/main.go

# Build the application
build:
	@echo "Building fluffy-mouton-api..."
	@go build -o bin/$(APP_NAME) $(MAIN_FILE)

# Run the application directly
start:
	@echo "Starting fluffy-mouton-api..."
	@go run $(MAIN_FILE)

# Run tests
test:
	@echo "Testing fluffy-mouton-api..."
	@go test -v ./...

# Clean the build directory
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf bin/

# Tidy module dependencies
tidy:
	@echo "Tidying..."
	@go mod tidy