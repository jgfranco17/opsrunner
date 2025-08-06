package executor

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type DefaultExecutor struct{}

func (c *DefaultExecutor) Run(ctx context.Context, name string, args ...string) (Result, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		// Get exit code if available
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			} else {
				// Unix-only, fallback
				exitCode = -1
			}
		} else {
			// Non-exit error, e.g., binary not found
			exitCode = -1
		}
	}

	return Result{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
	}, err
}
