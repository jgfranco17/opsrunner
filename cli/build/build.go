package build

import (
	"context"

	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
)

type Executor interface {
	Run(ctx context.Context, name string, args ...string) (executor.Result, error)
}

func Run(ctx context.Context, shellExecutor Executor, config *config.ProjectDefinition) (executor.Result, error) {
	return executor.Result{}, nil
}
