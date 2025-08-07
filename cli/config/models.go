package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gtithub.com/jgfranco17/opsrunner/cli/executor"
	"gtithub.com/jgfranco17/opsrunner/cli/logging"
	"gtithub.com/jgfranco17/opsrunner/cli/outputs"

	"gopkg.in/yaml.v3"
)

type ShellExecutor interface {
	Exec(ctx context.Context, command string) (executor.Result, error)
	AddEnv(env []string)
}

type ProjectDefinition struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Version     string   `yaml:"version"`
	RepoUrl     string   `yaml:"repo_url"`
	Codebase    Codebase `yaml:"codebase"`
}

// Load reads a YAML configuration from the provided reader and unmarshals
// it into a struct instance.
func Load(r io.Reader) (*ProjectDefinition, error) {
	var cfg ProjectDefinition
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}
	return &cfg, nil
}

type Codebase struct {
	Language     string    `yaml:"language"`
	Dependencies string    `yaml:"dependencies,omitempty"`
	Test         Operation `yaml:"test,omitempty"`
	Build        Operation `yaml:"build,omitempty"`
}

// Install installs the codebase dependencies based on the defined language.
func (c *Codebase) Install(ctx context.Context, executor ShellExecutor) error {
	logger := logging.FromContext(ctx)
	if c.Language == "" {
		return fmt.Errorf("codebase language is not defined")
	}
	switch strings.ToLower(c.Language) {
	case "go":
		logger.Debug("Running Go installation steps")
		result, err := executor.Exec(ctx, "go mod tidy")
		if err != nil {
			return fmt.Errorf("failed to install Go dependencies: %w", err)
		}
		result.PrintStdOut()
	case "python":
		logger.Debug("Running Python installation steps")
		pipCommand := fmt.Sprintf("pip install -r %s", c.Dependencies)
		result, err := executor.Exec(ctx, pipCommand)
		if err != nil {
			return fmt.Errorf("failed to install Go dependencies: %w", err)
		}
		result.PrintStdOut()
	default:
		return fmt.Errorf("unsupported codebase language: %s", c.Language)
	}
	return nil
}

type Operation struct {
	FailFast bool              `yaml:"fail_fast,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Steps    []string          `yaml:"steps"`
}

// Run executes the defined steps in the Operation using the provided envs.
func (op *Operation) Run(ctx context.Context, executor ShellExecutor) error {
	logger := logging.FromContext(ctx)

	env := os.Environ()
	if len(op.Env) > 0 {
		envsAdded := []string{}
		for k, v := range op.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
			envsAdded = append(envsAdded, k)
		}
		logger.Tracef("Loading additional %d additional environment variable(s): %v", len(op.Env), envsAdded)
	}
	executor.AddEnv(env)

	startTime := time.Now()
	var failedSteps []string
	for idx, step := range op.Steps {
		fmt.Printf("[%d] %s\n", idx+1, step)
		result, err := executor.Exec(ctx, step)
		if err != nil || result.ExitCode != 0 {
			if op.FailFast {
				return fmt.Errorf("error while running '%s' (exit code %d): %w", step, result.ExitCode, err)
			}
			failedSteps = append(failedSteps, step)
		}
		if result.Stdout != "" {
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", result.Stdout)
		}
		if result.Stderr != "" {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", result.Stderr)
		}
	}
	duration := time.Since(startTime)
	outputs.PrintTerminalWideLine("=")
	if len(failedSteps) > 0 {
		return fmt.Errorf("failed to run steps: %v", failedSteps)
	}
	outputs.PrintColoredMessage("green", "Ran %d tasks in %d ms", len(op.Steps), duration.Milliseconds())
	return nil
}
