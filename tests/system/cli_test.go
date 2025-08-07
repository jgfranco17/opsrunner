package cli_commands_test

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/core"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
)

var _ = Describe("CLI Commands System Tests", func() {
	var tempDir string
	var realExecutor *executor.DefaultExecutor
	var ctx context.Context

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		realExecutor = &executor.DefaultExecutor{}
		ctx = context.Background()
	})

	Describe("Build Command", func() {
		Context("when executing build command with valid configuration", func() {
			It("should execute build steps successfully", func() {
				// Given: A valid project configuration file
				projectConfig := &config.ProjectDefinition{
					Name:    "CLITestProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Install: config.Operation{
							Steps: []string{
								"echo 'Install step 1'",
								"echo 'Install step 2'",
							},
						},
						Build: config.Operation{
							Steps: []string{
								"echo 'Build step 1'",
								"echo 'Build step 2'",
								"echo 'Build completed'",
							},
						},
					},
				}

				// And: The configuration file is created
				content, err := yaml.Marshal(projectConfig)
				Expect(err).To(BeNil())
				configPath := filepath.Join(tempDir, ".opsrunner.yaml")
				err = os.WriteFile(configPath, content, 0644)
				Expect(err).To(BeNil())

				// When: The build command is executed
				buildOptions := &config.BuildOptions{
					NoInstall: false,
				}
				err = config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully
				Expect(err).To(BeNil())
			})
		})

		Context("when executing build command with no-install flag", func() {
			It("should skip installation steps", func() {
				// Given: A project configuration with both install and build steps
				projectConfig := &config.ProjectDefinition{
					Name:    "NoInstallProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Install: config.Operation{
							Steps: []string{
								"echo 'This should be skipped'",
								"exit 1", // This would fail if executed
							},
						},
						Build: config.Operation{
							Steps: []string{
								"echo 'Build step executed'",
								"echo 'Build completed'",
							},
						},
					},
				}

				// When: The build command is executed with no-install flag
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully without executing install steps
				Expect(err).To(BeNil())
			})
		})

		Context("when executing build command with failing steps", func() {
			It("should handle failures appropriately", func() {
				// Given: A project configuration with a failing step
				projectConfig := &config.ProjectDefinition{
					Name:    "FailingProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Build: config.Operation{
							FailFast: true,
							Steps: []string{
								"echo 'Step 1'",
								"exit 1",        // This will fail
								"echo 'Step 3'", // This should not execute
							},
						},
					},
				}

				// When: The build command is executed
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The build should fail
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to run build steps"))
			})
		})
	})

	Describe("Command Registry", func() {
		Context("when creating a command registry", func() {
			It("should register commands correctly", func() {
				// Given: A command registry with project information
				registry := core.NewCommandRegistry("test-project", "Test project description", "1.0.0")

				// When: Commands are registered
				buildCommand := core.GetBuildCommand(realExecutor)
				commands := []*cobra.Command{buildCommand}
				registry.RegisterCommands(commands)

				// Then: The registry should contain the registered commands
				Expect(registry).ToNot(BeNil())
				Expect(buildCommand).ToNot(BeNil())
				Expect(buildCommand.Use).To(Equal("build"))
			})
		})
	})
})
