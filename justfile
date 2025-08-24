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

# Run all BDD tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Run specific test category
bdd category="acceptance":
    #!/usr/bin/env bash
    echo "Running {{ category }} tests..."
    if ! command -v ginkgo &> /dev/null; then
        echo "Ginkgo CLI not found."
        echo "Install with: go install github.com/onsi/ginkgo/v2/ginkgo@latest"
        echo "After installation, ensure that GOPATH is on your PATH."
    fi
    case "{{ category }}" in
        "integration")
            ginkgo -r -v ./tests/integration/...
            ;;
        "system")
            ginkgo -r -v ./tests/system/...
            ;;
        "acceptance")
            ginkgo -r -v ./tests/acceptance/...
            ;;
        *)
            echo "Unknown test category: {{ category }}"
            echo "Available categories: unit, integration, system, acceptance"
            exit 1
            ;;
    esac

# Run tests with coverage
test-coverage:
    #!/usr/bin/env bash
    go test -coverpkg="./..." -coverprofile="coverage.out" -covermode="count" ./...
    go tool cover -html="coverage.out" -o coverage.html
    xdg-open coverage.html

# Run tests in parallel
bdd-parallel:
    #!/usr/bin/env bash
    echo "Running tests in parallel..."
    ginkgo -r -v -p ./tests/...

# Build CLI binary
build-bin version="0.0.0-dev":
    #!/usr/bin/env bash
    BIN_NAME="opsrunner"
    echo "Building {{ PROJECT_NAME }} binary..."
    go mod download all
    VERSION=$(jq -r .version specs.json)
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o "./${BIN_NAME}" main.go
    chmod +x "./${BIN_NAME}"
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
