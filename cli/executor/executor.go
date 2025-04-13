package executor

import (
	"context"
	"os/exec"
)

type CommandExecutor struct{}

func (c CommandExecutor) Exec(ctx context.Context, name string, args string) (int, string, error) {
	cmd := exec.CommandContext(ctx, name, args)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), string(output), exitErr
		}
		return -1, "", err
	}
	return 0, string(output), nil
}
