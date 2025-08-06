package config

import (
	"context"
	"fmt"
	"io"
	"time"

	"gtithub.com/jgfranco17/opsrunner/cli/outputs"

	"gopkg.in/yaml.v3"
)

type ShellExecutor interface {
	Exec(ctx context.Context, name string, args string) (int, string, error)
}

type ProjectDefinition struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Version     string  `yaml:"version"`
	RepoUrl     string  `yaml:"repo_url"`
	Test        Test    `yaml:"test,omitempty"`
	Build       Build   `yaml:"build,omitempty"`
	Package     Package `yaml:"package,omitempty"`
}

func Load(r io.Reader) (*ProjectDefinition, error) {
	var cfg ProjectDefinition
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}
	return &cfg, nil
}

type Build struct {
}

type Test struct {
}

type Package struct {
}

type Step struct {
	Command string `yaml:"cmd"`
	Args    string `yaml:"args,omitempty"`
}

type Operation struct {
	Env   map[string]string `yaml:"env,omitempty"`
	Steps []Step            `yaml:"steps"`
}

func (t *Operation) Run(ctx context.Context, executor ShellExecutor) error {
	if len(t.Env) > 0 {
		restoreFunc, err := WithTempEnv(ctx, t.Env)
		if err != nil {
			return fmt.Errorf("failed to set temporary env: %w", err)
		}
		defer restoreFunc()
	}
	outputs.PrintTerminalWideLine("=")

	startTime := time.Now()
	for idx, step := range t.Steps {
		fmt.Printf("[%d] %s %s\n", idx+1, step.Command, step.Args)
		exitCode, output, err := executor.Exec(ctx, step.Command, step.Args)
		if err != nil || exitCode != 0 {
			return fmt.Errorf("error while running '%s' (exit code %d): %w", step.Command, exitCode, err)
		}
		fmt.Println(string(output))
	}
	duration := time.Since(startTime)
	outputs.PrintTerminalWideLine("=")
	outputs.PrintColoredMessage("green", "OK", "Ran %d tasks in %d ms", len(t.Steps), duration.Milliseconds())
	return nil
}
