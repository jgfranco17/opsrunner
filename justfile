PROJECT_NAME := "opsrunner"

# Default command
_default:
    @just --list --unsorted

# Sync Go modules
tidy:
    go mod tidy
    @echo "All modules synced, Go workspace ready!"

# CLI local run wrapper
cli *args:
    @go run main.go {{ args }}

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Run coverage and open a report
view-coverage:
    go clean -testcache
    go test -coverpkg="./..." -coverprofile="coverage.out" -covermode="count" ./...
    go tool cover -html="coverage.out" -o coverage.html
    xdg-open coverage.html

# Build CLI binary
build-bin:
    #!/usr/bin/env bash
    echo "Building {{ PROJECT_NAME }} binary..."
    go mod download all
    VERSION=$(jq -r .version specs.json)
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o ./{{ PROJECT_NAME }} main.go
    echo "Built binary for {{ PROJECT_NAME }} ${VERSION} successfully!"

# Build the Docker image
docker-build:
    docker build --build-arg VERSION=0.0.0-img --no-cache -t opsrunner:dev .
    @echo "Built OpsRunner image successfully!"

# Generate CLI documentation
generate-docs directory="docs":
    #!/usr/bin/env bash
    just cli docs {{ directory }}
    find {{ directory }} -name '*completion*' -exec rm -f {} \;
