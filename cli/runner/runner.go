package runner

import (
	"context"
	"fmt"
	"os"
	"time"

	"gtithub.com/jgfranco17/opsrunner/cli/outputs"

	"gopkg.in/yaml.v3"
)

type ShellExecutor interface {
	Exec(ctx context.Context, name string, args string) (int, string, error)
}

type OpsRunnerConfig struct {
	RepoUrl string          `yaml:"repo_url"`
	Tasks   map[string]Task `yaml:"tasks"`
}

func ReadConfigFromFile(ctx context.Context, filePath string) (*OpsRunnerConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read YAML file: %w", err)
	}
	var result OpsRunnerConfig
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML: %w", err)
	}
	return &result, nil
}

type Step struct {
	Command string `yaml:"cmd"`
	Args    string `yaml:"args,omitempty"`
}

type Task struct {
	Description string            `yaml:"description,omitempty"`
	Category    string            `yaml:"category,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Steps       []Step            `yaml:"steps"`
}

func (t *Task) Run(ctx context.Context, executor ShellExecutor) error {
	if len(t.Env) > 0 {
		restoreFunc, err := WithTempEnv(ctx, t.Env)
		if err != nil {
			return fmt.Errorf("Failed to set temporary env: %w", err)
		}
		defer restoreFunc()
	}
	outputs.PrintTerminalWideLine("=")

	startTime := time.Now()
	for idx, step := range t.Steps {
		fmt.Printf("[%d] %s %s\n", idx+1, step.Command, step.Args)
		exitCode, output, err := executor.Exec(ctx, step.Command, step.Args)
		if err != nil || exitCode != 0 {
			return fmt.Errorf("Error while running '%s' (exit code %d): %w", step.Command, exitCode, err)
		}
		fmt.Println(string(output))
	}
	duration := time.Since(startTime)
	outputs.PrintTerminalWideLine("=")
	outputs.PrintColoredMessage("green", "OK", "Ran %d tasks in %d ms", len(t.Steps), duration.Milliseconds())
	return nil
}
