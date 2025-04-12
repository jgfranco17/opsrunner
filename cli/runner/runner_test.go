package runner

import (
	"context"
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
	step := Step{
		Command: "echo",
		Args:    "hello",
	}
	exitCode, output, err := step.Exec()
	assert.Equal(t, 0, exitCode)
	assert.NoError(t, err)
	assert.Contains(t, output, "hello")
}

func TestStepFailOk(t *testing.T) {
	step := Step{
		Command: "exit",
		Args:    "",
	}
	exitCode, _, err := step.Exec()
	assert.Equal(t, -1, exitCode)
	assert.Error(t, err)
}
