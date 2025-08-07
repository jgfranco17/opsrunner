package build

import (
	"context"
	"fmt"

	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
	"gtithub.com/jgfranco17/opsrunner/cli/logging"
)

type Executor interface {
	Exec(ctx context.Context, command string) (executor.Result, error)
}

type Options struct {
	NoInstall bool
}

func Exec(ctx context.Context, shellExecutor config.ShellExecutor, config *config.ProjectDefinition, opts *Options) error {
	logger := logging.FromContext(ctx)

	if opts == nil {
		opts = &Options{}
	}
	if !opts.NoInstall {
		logger.Debug("Installing codebase dependencies")
		if err := config.Codebase.Install(ctx, shellExecutor); err != nil {
			return fmt.Errorf("failed to install codebase dependencies: %w", err)
		}
	} else {
		logger.Info("Skipping codebase dependency installation")
	}
	if len(config.Codebase.Build.Steps) == 0 {
		logger.Warn("No build steps defined in the configuration.")
	}
	err := config.Codebase.Build.Run(ctx, shellExecutor)
	if err != nil {
		return fmt.Errorf("failed to run build steps: %w", err)
	}
	return nil
}
