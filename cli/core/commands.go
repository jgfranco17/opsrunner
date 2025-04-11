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
		Short: "Ping a target URL",
		Long:  "Check if a target URL is live and responds with a 2xx status code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Not enough arguments, expected 1 but got %d", len(args))
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
			err = taskToRun.Run()
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
