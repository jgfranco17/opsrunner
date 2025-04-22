package core

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"gtithub.com/jgfranco17/opsrunner/cli/runner"
)

type ShellExecutor interface {
	Exec(ctx context.Context, name string, args string) (int, string, error)
}

func GetRunCommand(shellExecutor ShellExecutor) *cobra.Command {
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
			err = taskToRun.Run(ctx, shellExecutor)
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

func GetGenerateDocsCommand(rootCmd *cobra.Command) *cobra.Command {
	var outputDirPath string
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate CLI documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := GenerateDocs(rootCmd, outputDirPath); err != nil {
				return fmt.Errorf("could not generate CLI docs: %w", err)
			}
			return nil
		},
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVarP(&outputDirPath, "output", "o", "docs", "Output directory")
	return cmd
}
