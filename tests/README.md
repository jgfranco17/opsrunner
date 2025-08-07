# OpsRunner BDD Testing Documentation

## Overview

This document outlines the Behavior-Driven Development (BDD) testing strategy implemented for the OpsRunner CLI system using Ginkgo and Gomega frameworks.

## What is BDD Testing?

Behavior-Driven Development (BDD) is a software development methodology that encourages collaboration between developers, QA engineers, and non-technical stakeholders. BDD focuses on describing software behavior in natural language that both technical and non-technical team members can understand.

### Key Benefits of BDD Testing

1. **Improved Communication**: BDD tests are written in natural language (Given-When-Then format), making them accessible to all stakeholders
2. **Living Documentation**: Tests serve as executable specifications that are always up-to-date
3. **Early Bug Detection**: BDD tests help identify issues early in the development cycle
4. **Regression Prevention**: Automated BDD tests prevent regressions when new features are added
5. **Better Test Coverage**: BDD encourages testing from a user's perspective, leading to more comprehensive coverage

## Testing Strategy for OpsRunner

### Test Categories

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

## Test Structure

### Given-When-Then Format

All BDD tests follow the Given-When-Then format:

```go
Describe("Build Command", func() {
    Context("when a valid configuration file is provided", func() {
        It("should execute build steps successfully", func() {
            // Given: A valid YAML configuration file
            configFile := createValidConfigFile()

            // When: The build command is executed
            result := executeBuildCommand(configFile)

            // Then: The build should complete successfully
            Expect(result.ExitCode).To(Equal(0))
            Expect(result.Stdout).To(ContainSubstring("Build completed successfully"))
        })
    })
})
```

### Test Organization

```text
tests/
├── unit/                    # Unit tests for individual components
│   ├── config/             # Configuration loading and validation tests
│   ├── executor/           # Executor component tests
│   └── commands/           # Command structure tests
├── integration/            # Integration tests
│   ├── build_workflow/     # End-to-end build workflow tests
│   └── config_executor/    # Configuration and executor integration
├── system/                 # System-level tests
│   ├── cli_commands/       # Full CLI command tests
│   └── file_operations/    # File system operation tests
├── acceptance/             # Acceptance tests
│   ├── user_stories/       # User story validation
│   └── business_flows/     # Business workflow tests
├── fixtures/               # Test data and fixtures
│   ├── configs/           # Sample configuration files
│   └── expected_outputs/  # Expected test outputs
└── helpers/               # Test helper functions
    ├── test_utils.go      # Common test utilities
    └── mock_executor.go   # Mock executor implementation
```

## Why These Tests Are Important

### 1. **Reliability Assurance**

- **Problem**: CLI tools can fail in production due to unexpected inputs or environment issues
- **Solution**: BDD tests validate behavior under various conditions
- **Benefit**: Reduces production failures and improves system reliability

### 2. **Regression Prevention**

- **Problem**: New features can break existing functionality
- **Solution**: Comprehensive test suite catches regressions early
- **Benefit**: Maintains system stability during development

### 3. **Documentation and Onboarding**

- **Problem**: New team members struggle to understand system behavior
- **Solution**: BDD tests serve as living documentation
- **Benefit**: Faster onboarding and better knowledge transfer

### 4. **Quality Gates**

- **Problem**: Code changes can introduce bugs without proper validation
- **Solution**: Automated tests run on every commit
- **Benefit**: Ensures code quality and prevents broken builds

### 5. **User Experience Validation**

- **Problem**: Technical implementation may not match user expectations
- **Solution**: BDD tests focus on user behavior and outcomes
- **Benefit**: Ensures the system meets user needs

## Test Execution

### Running All Tests

```bash
just test
```

### Running Specific Test Categories
```bash
# Unit tests only
just test-category unit

# Integration tests only
just test-category integration

# System tests only
just test-category system

# Acceptance tests only
just test-category acceptance
```

### Running Tests with Advanced Options
```bash
# Run tests with coverage
just test-coverage

# Run tests in parallel
just bdd-parallel

# Run tests with verbose output
just test-verbose

# Run tests and generate JUnit report
just test-junit

# Run tests with focus pattern
just test-focus 'Build'

# Skip tests matching pattern
just test-skip 'System'
```

### Setup and Help
```bash
# Setup test environment
just test-setup

# Show test help
just test-help

# Install Ginkgo CLI
just install-ginkgo
```

## Continuous Integration

### GitHub Actions Integration
Tests are automatically run on:
- Every pull request
- Every push to main branch
- Scheduled nightly runs

### Quality Gates
- All tests must pass before merging
- Minimum code coverage requirements
- Performance benchmarks must be met

## Best Practices

### 1. **Test Isolation**
- Each test should be independent
- Use `t.TempDir()` for temporary directories (automatic cleanup)
- Avoid shared state between tests

### 2. **Descriptive Test Names**
- Use clear, descriptive test names
- Follow the Given-When-Then pattern
- Make test intentions obvious

### 3. **Test Data Management**
- Use fixtures for test data
- Keep test data minimal and focused
- Document test data purpose

### 4. **Mock Usage**
- Mock external dependencies
- Use realistic mock behavior
- Document mock expectations

### 5. **Error Testing**
- Test both success and failure scenarios
- Validate error messages and codes
- Test edge cases and boundary conditions

## Metrics and Reporting

### Test Coverage
- Aim for >90% code coverage
- Focus on critical path coverage
- Regular coverage reporting

### Test Performance
- Monitor test execution time
- Optimize slow tests
- Parallel test execution where possible

### Quality Metrics
- Test failure rate
- Test maintenance effort
- Bug detection effectiveness

## Conclusion

BDD testing with Ginkgo provides a robust foundation for ensuring OpsRunner's reliability, maintainability, and user satisfaction. By focusing on behavior rather than implementation details, these tests serve as both quality assurance tools and living documentation that evolves with the system.

The comprehensive test suite helps maintain high code quality, prevents regressions, and ensures that OpsRunner continues to meet user expectations as it evolves.
