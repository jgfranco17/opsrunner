PROJECT_NAME := "opsrunner"

# Default command
_default:
    @just --list --unsorted

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./cli/...

# Sync Go modules
tidy:
    cd cli && go mod tidy
    go mod tidy
    go work sync

# CLI local run wrapper
cli *args:
    @go run main.go {{ args }}

# Build CLI binary
build:
    #!/usr/bin/env bash
    echo "Building {{ PROJECT_NAME }} binary..."
    go mod download all
    VERSION=$(jq -r .version specs.json)
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o ./{{ PROJECT_NAME }} main.go
    echo "Built binary for {{ PROJECT_NAME }} ${VERSION} successfully!"
