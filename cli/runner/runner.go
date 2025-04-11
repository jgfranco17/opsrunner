package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

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

type Task struct {
	Description string            `yaml:"description,omitempty"`
	Category    string            `yaml:"category,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Steps       []Step            `yaml:"steps"`
}

func (t *Task) Run() error {
	for envVarible, value := range t.Env {
		fmt.Printf("%s=%s\n", envVarible, value)
	}

	for idx, step := range t.Steps {
		fmt.Printf("[%d / %d] Running: %s %s\n", idx+1, len(t.Steps), step.Command, step.Args)
		cmdExec := exec.Command(step.Command, step.Args)
		output, err := cmdExec.Output()
		if err != nil {
			return fmt.Errorf("Error while running '%s': %w", step.Command, err)
		}
		fmt.Println(string(output))
	}
	return nil
}

type Step struct {
	Command string `yaml:"cmd"`
	Args    string `yaml:"args,omitempty"`
}
