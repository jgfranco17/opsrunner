package core

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"cli/runner"
)

func GetRunCommnd() *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a task",
		Long:  "Read the config file and run a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Too many arguments, expected 1 but got %d", len(args))
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			config, err := runner.ReadConfigFromFile(ctx, filePath)
			if err != nil {
				return err
			}
			taskName := args[0]
			taskToRun, ok := config.Tasks[taskName]
			if !ok {
				return fmt.Errorf("No such task: %s", taskName)
			}
			err = taskToRun.Run(ctx)
			if err != nil {
				return err
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "opsrunner.yaml", "OpsRunner definition file")
	return cmd
}
