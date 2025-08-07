package core

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gtithub.com/jgfranco17/opsrunner/cli/build"
	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
	"gtithub.com/jgfranco17/opsrunner/cli/logging"
)

type BashExecutor interface {
	Exec(ctx context.Context, command string) (executor.Result, error)
	AddEnv(env []string)
}

func GetBuildCommand(shellExecutor BashExecutor) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Run the build operations",
		Long:  "Read the config file and run the build operations defined in it.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.FromContext(cmd.Context())
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			logger.Debugf("Starting build with config file: %s", filePath)
			contents, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", filePath, err)
			}
			config, err := config.Load(contents)
			if err != nil {
				return fmt.Errorf("failed to load config from file: %w", err)
			}
			if err := build.Exec(ctx, shellExecutor, config); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", ".opsrunner.yaml", "OpsRunner definition file")
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
