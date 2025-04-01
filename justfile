# Roguelike Game Justfile
# Run commands with 'just <command>'

# Default recipe to run when just is called without arguments
default:
    @just --list

# Build the game
build:
    cd roguelike-gruid-project && go build -o ../roguelike

# Run the game
run: build
    ./roguelike

# Build and run with race detection
run-race:
    cd roguelike-gruid-project && go run -race .

# Clean build artifacts
clean:
    rm -f roguelike

# Run tests
test:
    cd roguelike-gruid-project && go test ./...

# Run tests with verbose output
test-verbose:
    cd roguelike-gruid-project && go test -v ./...

# Format code
fmt:
    cd roguelike-gruid-project && go fmt .
    @echo "Go files formatted successfully"

# Check for code issues
lint:
    @echo "Running go vet..."
    cd roguelike-gruid-project && go vet .
    @echo "Running go vet completed successfully"

# Install dependencies
deps:
    cd roguelike-gruid-project && go mod tidy

# Build for web (WebAssembly)
build-wasm:
    cd roguelike-gruid-project && GOOS=js GOARCH=wasm go build -o ../roguelike.wasm
