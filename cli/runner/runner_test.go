package runner

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigOk(t *testing.T) {
	config, err := ReadConfigFromFile(context.Background(), "./resources/simple.yaml")
	assert.NoError(t, err, "Unexpected error while loading config from file")
	assert.Equal(t, "www.example.com", config.RepoUrl)
	assert.Equal(t, 1, len(config.Tasks))
}

func TestLoadConfigFail_InvalidPath(t *testing.T) {
	config, err := ReadConfigFromFile(context.Background(), "non-existent-file")
	assert.ErrorContains(t, err, "no such file or directory")
	assert.Empty(t, config)
}

func TestLoadConfigFail_InvalidYamlSchema(t *testing.T) {
	config, err := ReadConfigFromFile(context.Background(), "./resources/invalid.txt")
	assert.ErrorContains(t, err, "Failed to unmarshal YAML")
	assert.Empty(t, config)
}

func TestStepExecOk(t *testing.T) {
	task := Task{
		Description: "A simple test task",
		Env:         make(map[string]string),
		Category:    "",
		Steps: []Step{
			{
				Command: "some-command",
				Args:    "--arg value",
			},
		},
	}
	mockExecutor := NewMockExecutor(1).WithStep("some-command", "--arg value", 0, "Ran some-command!", nil)
	err := task.Run(context.Background(), mockExecutor)
	assert.NoError(t, err)
	assert.True(t, mockExecutor.IsCalled())
}

func TestStepExecWithStepError(t *testing.T) {
	task := Task{
		Description: "A simple test task",
		Env:         make(map[string]string),
		Category:    "",
		Steps: []Step{
			{
				Command: "some-command",
				Args:    "--arg value",
			},
		},
	}
	mockExecutor := NewMockExecutor(1).WithStep("some-command", "--arg value", 1, "Something went wrong!", fmt.Errorf("Some error"))
	err := task.Run(context.Background(), mockExecutor)
	assert.True(t, mockExecutor.IsCalled())
	assert.ErrorContains(t, err, "Error while running 'some-command' (exit code 1)")
}
