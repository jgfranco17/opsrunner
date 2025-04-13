package core

import (
	"bytes"
	"cli/runner"
	"testing"

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
