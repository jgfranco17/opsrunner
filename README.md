# OpsRunner CLI

_DevOps on local, simplified._

![STATUS](https://img.shields.io/badge/status-active-brightgreen?style=for-the-badge)
![LICENSE](https://img.shields.io/badge/license-BSD3-blue?style=for-the-badge)

---

## Introduction

### About

OpsRunner is a powerful, user-friendly Command-Line Interface (CLI) tool tailored for
CI/CD pipelines and build management. Designed to simplify the orchestration of complex
build processes, OpsRunner offers a structured and reliable alternative to traditional
scripting. By leveraging project specifications defined in YAML files, this tool empowers
developers to streamline repetitive tasks and enhance automation across various
environments.

Whether you're running a local build or deploying to production, OpsRunner makes it easier
to achieve consistency, traceability, and scalability in your workflows.

## Installation

To download the CLI, an install script has been provided.

```bash
wget -O - https://raw.githubusercontent.com/jgfranco17/opsrunner/main/install.sh | bash
```

They always say not to just blindly run scripts from the internet, so feel free to examine
the file first before running.

> [!NOTE]
> This CLI is still an alpha prototype.

## Testing

### Test Categories

While smaller units are tested with simple Go Test assertions, we apply BDD approach for
the larger interactions. BDD testing with Ginkgo provides a robust foundation for ensuring
OpsRunner's reliability, maintainability, and user satisfaction. By focusing on behavior
rather than implementation details, these tests serve as both quality assurance tools and
living documentation that evolves with the system.

#### 1. Unit Tests (Component Level)

- **Purpose**: Test individual components in isolation
- **Scope**: Functions, methods, and small units of code
- **Tools**: Go Test
- **Location**: Alongside source code in `cli/`

#### 2. Integration Tests (Component Interaction)

- **Purpose**: Test how components work together
- **Scope**: Command execution, configuration loading, executor interactions
- **Tools**: Ginkgo + Gomega + Mock objects
- **Location**: Alongside source code in `cli/`

#### 3. System Tests (End-to-End)

- **Purpose**: Test the complete system from a user's perspective
- **Scope**: Full CLI commands, file operations, real shell execution
- **Tools**: Ginkgo + Gomega + Real shell executor
- **Location**: `tests/system/`

#### 4. Acceptance Tests (Business Requirements)

- **Purpose**: Verify that the system meets business requirements
- **Scope**: User stories, feature requirements, business workflows
- **Tools**: Ginkgo + Gomega + Test scenarios
- **Location**: `tests/acceptance/`

### Execution

To run the tests, you can use the `just` command. The `justfile` contains various targets
for different test categories.

```bash
# Run standard assertions with go-test
just test

# Run a specific category of BDD tests
just bdd category="integration"  # For integration tests
just bdd category="system"  # For system tests
just bdd category="acceptance"  # For acceptance tests

# Run all BDD tests in parallel
just bdd-parallel
```

### Automation

#### GitHub Actions Integration

Tests are automatically run on:

- Every pull request
- Every push to main branch
- Scheduled nightly runs

#### Quality Gates

- All tests must pass before merging
- Minimum code coverage requirements
- Performance benchmarks must be met

## License

This project is licensed under the BSD-3 License. See the LICENSE file for more details.
