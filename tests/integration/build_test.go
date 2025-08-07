package build_workflow_test

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
	"gtithub.com/jgfranco17/opsrunner/tests/internal/helpers"
)

var _ = Describe("Build Workflow Integration", func() {
	var tempDir string
	var mockExecutor *helpers.MockExecutor
	var ctx context.Context

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		mockExecutor = helpers.CreateMockExecutor()
		ctx = context.Background()
	})

	Describe("Complete Build Process", func() {
		Context("when a valid project configuration is provided", func() {
			It("should execute the complete build workflow successfully", func() {
				// Given: A complete project configuration with install and build steps
				projectConfig := &config.ProjectDefinition{
					Name:        "TestProject",
					Version:     "1.0.0",
					Description: "A test project for integration testing",
					RepoUrl:     "https://github.com/test/project",
					Codebase: config.Codebase{
						Language:     "go",
						Dependencies: "go.mod",
						Install: config.Operation{
							FailFast: true,
							Env: map[string]string{
								"GO_ENV": "test",
							},
							Steps: []string{
								"echo 'Installing dependencies...'",
								"go mod tidy",
								"go mod download",
							},
						},
						Build: config.Operation{
							FailFast: false,
							Env: map[string]string{
								"BUILD_ENV":   "production",
								"CGO_ENABLED": "0",
							},
							Steps: []string{
								"echo 'Building project...'",
								"echo 'go build -o testapp'",
								"echo 'Build completed successfully'",
							},
						},
					},
				}

				// And: A configuration file is created
				content, err := yaml.Marshal(projectConfig)
				Expect(err).To(BeNil())
				configPath := filepath.Join(tempDir, ".opsrunner.yaml")
				err = os.WriteFile(configPath, content, 0644)
				Expect(err).To(BeNil())
				_, err = os.Stat(configPath)
				Expect(err).To(BeNil())

				// When: The build process is executed
				buildOptions := &config.BuildOptions{
					NoInstall: false,
				}
				err = config.Build(ctx, mockExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully
				Expect(err).To(BeNil())

				// And: All expected commands should be executed
				executions := mockExecutor.GetExecutions()
				Expect(executions).To(HaveLen(6)) // 3 install + 3 build steps

				// Verify install steps
				Expect(executions[0].Command).To(Equal("echo 'Installing dependencies...'"))
				Expect(executions[1].Command).To(Equal("go mod tidy"))
				Expect(executions[2].Command).To(Equal("go mod download"))

				// Verify build steps
				Expect(executions[3].Command).To(Equal("echo 'Building project...'"))
				Expect(executions[4].Command).To(Equal("echo 'go build -o testapp'"))
				Expect(executions[5].Command).To(Equal("echo 'Build completed successfully'"))

				// And: Environment variables should be properly set
				env := mockExecutor.GetEnv()
				Expect(env).To(ContainElement("GO_ENV=test"))
				Expect(env).To(ContainElement("BUILD_ENV=production"))
				Expect(env).To(ContainElement("CGO_ENABLED=0"))
			})
		})

		Context("when build is executed with no-install flag", func() {
			It("should skip installation and only execute build steps", func() {
				// Given: A project configuration with both install and build steps
				projectConfig := &config.ProjectDefinition{
					Name:    "SkipInstallProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Install: config.Operation{
							Steps: []string{
								"echo 'This should be skipped'",
								"go mod tidy",
							},
						},
						Build: config.Operation{
							Steps: []string{
								"echo 'Building without install'",
								"echo 'go build'",
							},
						},
					},
				}

				// When: The build process is executed with no-install flag
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, mockExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully
				Expect(err).To(BeNil())

				// And: Only build steps should be executed
				executions := mockExecutor.GetExecutions()
				Expect(executions).To(HaveLen(2)) // Only 2 build steps

				// Verify only build steps were executed
				Expect(executions[0].Command).To(Equal("echo 'Building without install'"))
				Expect(executions[1].Command).To(Equal("echo 'go build'"))

				// And: Install steps should not be executed
				for _, execution := range executions {
					Expect(execution.Command).ToNot(Equal("echo 'This should be skipped'"))
					Expect(execution.Command).ToNot(Equal("go mod tidy"))
				}
			})
		})

		Context("when a build step fails with fail_fast enabled", func() {
			It("should stop execution and return an error", func() {
				// Given: A project configuration with a failing build step
				projectConfig := &config.ProjectDefinition{
					Name:    "FailingBuildProject",
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
							FailFast: true,
							Steps: []string{
								"echo 'Build step 1'",
								"exit 1",              // This will fail
								"echo 'Build step 3'", // This should not execute
							},
						},
					},
				}

				// Mock the executor to return an error for the failing step
				mockExecutor.SetExecFunc(func(ctx context.Context, command string) (executor.Result, error) {
					if command == "exit 1" {
						result := executor.Result{
							Stdout:   "",
							Stderr:   "command failed",
							ExitCode: 1,
						}
						mockExecutor.RecordExecution(command, result, nil)
						return result, nil
					}

					result := executor.Result{
						Stdout:   "success",
						Stderr:   "",
						ExitCode: 0,
					}
					mockExecutor.RecordExecution(command, result, nil)
					return result, nil
				})

				// When: The build process is executed
				buildOptions := &config.BuildOptions{
					NoInstall: false,
				}
				err := config.Build(ctx, mockExecutor, projectConfig, buildOptions)

				// Then: The build should fail
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to run build steps"))

				// And: Execution should stop at the failing step
				executions := mockExecutor.GetExecutions()
				Expect(executions).To(HaveLen(4)) // 2 install + 2 build steps (stopped at failure)

				// Verify install steps completed
				Expect(executions[0].Command).To(Equal("echo 'Install step 1'"))
				Expect(executions[1].Command).To(Equal("echo 'Install step 2'"))

				// Verify build steps stopped at failure
				Expect(executions[2].Command).To(Equal("echo 'Build step 1'"))
				Expect(executions[3].Command).To(Equal("exit 1"))

				// And: The third build step should not be executed
				for _, execution := range executions {
					Expect(execution.Command).ToNot(Equal("echo 'Build step 3'"))
				}
			})
		})
	})
})
