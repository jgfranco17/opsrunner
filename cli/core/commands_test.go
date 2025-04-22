package core

import (
	"bytes"
	"testing"

	"gtithub.com/jgfranco17/opsrunner/cli/runner"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type CliCommandFunction func() *cobra.Command

type CommandRunner func(cmd *cobra.Command, args []string)

type CliRunResult struct {
	ShellOutput string
	Error       error
}

// Helper function to simulate CLI execution
func ExecuteTestCommand(t *testing.T, cmd *cobra.Command, args ...string) CliRunResult {
	t.Helper()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	_, err := cmd.ExecuteC()
	return CliRunResult{
		ShellOutput: buf.String(),
		Error:       err,
	}
}

func TestRunCommandDefaultSuccess(t *testing.T) {
	mockExecutor := runner.NewMockExecutor(1).WithStep("some-command", "--arg value", 0, "Ran some-command!", nil)
	result := ExecuteTestCommand(t, GetRunCommand(mockExecutor), "test", "-f", "./resources/simple.yaml")
	assert.NoError(t, result.Error, "Unexpected error while executing run command")
}

func TestRunCommandFail_TooManyArguments(t *testing.T) {
	mockExecutor := runner.NewMockExecutor(0)
	result := ExecuteTestCommand(t, GetRunCommand(mockExecutor), "test", "another", "test")
	assert.ErrorContains(t, result.Error, "accepts 1 arg(s), received 3")
}

func TestRunCommandFail_InvalidTaskName(t *testing.T) {
	mockExecutor := runner.NewMockExecutor(0)
	result := ExecuteTestCommand(t, GetRunCommand(mockExecutor), "non-existent", "-f", "./resources/simple.yaml")
	assert.ErrorContains(t, result.Error, "No such task: non-existent")
}
