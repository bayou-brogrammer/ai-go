# Roguelike Game Justfile
# Run commands with 'just <command>'

# Default recipe to run when just is called without arguments
default:
    @just --list

# Build the game
build:
    go build -o roguelike

# Run the game
run: build
    ./roguelike

# Build and run with race detection
run-race:
    go run -race .

# Clean build artifacts
clean:
    rm -f roguelike

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Format code
fmt:
    go fmt ./...

# Check for code issues
lint:
    go vet ./...

# Install dependencies
deps:
    go mod tidy

# Build for web (WebAssembly)
build-wasm:
    GOOS=js GOARCH=wasm go build -o roguelike.wasm
