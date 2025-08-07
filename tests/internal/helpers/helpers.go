package helpers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
)

// CreateTestConfigFile creates a test configuration file using the existing config structure
func CreateTestConfigFile(t *testing.T, dir string, projectConfig *config.ProjectDefinition) string {
	t.Helper()

	content, err := yaml.Marshal(projectConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	configPath := filepath.Join(dir, ".opsrunner.yaml")
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	return configPath
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("Condition not met within %v: %s", timeout, message)
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("Expected file to not exist: %s", path)
	}
}

// ReadFileContent reads the content of a file
func ReadFileContent(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

// CreateMockExecutor creates a mock executor for testing
func CreateMockExecutor() *MockExecutor {
	return &MockExecutor{
		executions: make([]ExecutionRecord, 0),
		env:        make([]string, 0),
	}
}

// MockExecutor is a mock implementation of the executor interface for testing
type MockExecutor struct {
	executions []ExecutionRecord
	env        []string
	execFunc   func(ctx context.Context, command string) (executor.Result, error)
}

// ExecutionRecord records an execution for verification
type ExecutionRecord struct {
	Command string
	Result  executor.Result
	Error   error
}

// Exec implements the executor interface for mocking
func (m *MockExecutor) Exec(ctx context.Context, command string) (executor.Result, error) {
	// Use custom exec function if provided
	if m.execFunc != nil {
		return m.execFunc(ctx, command)
	}

	// Default successful result
	result := executor.Result{
		Stdout:   "mock stdout",
		Stderr:   "",
		ExitCode: 0,
	}

	// Record the execution
	record := ExecutionRecord{
		Command: command,
		Result:  result,
		Error:   nil,
	}
	m.executions = append(m.executions, record)

	return result, nil
}

// AddEnv adds environment variables to the mock executor
func (m *MockExecutor) AddEnv(env []string) {
	m.env = append(m.env, env...)
}

// GetExecutions returns all recorded executions
func (m *MockExecutor) GetExecutions() []ExecutionRecord {
	return m.executions
}

// GetLastExecution returns the last recorded execution
func (m *MockExecutor) GetLastExecution() *ExecutionRecord {
	if len(m.executions) == 0 {
		return nil
	}
	return &m.executions[len(m.executions)-1]
}

// ClearExecutions clears all recorded executions
func (m *MockExecutor) ClearExecutions() {
	m.executions = make([]ExecutionRecord, 0)
}

// SetExecFunc sets a custom execution function for the mock
func (m *MockExecutor) SetExecFunc(fn func(ctx context.Context, command string) (executor.Result, error)) {
	m.execFunc = fn
}

// GetEnv returns the environment variables
func (m *MockExecutor) GetEnv() []string {
	return m.env
}

// RecordExecution manually records an execution for testing
func (m *MockExecutor) RecordExecution(command string, result executor.Result, err error) {
	record := ExecutionRecord{
		Command: command,
		Result:  result,
		Error:   err,
	}
	m.executions = append(m.executions, record)
}
