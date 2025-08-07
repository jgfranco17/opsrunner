package user_stories_test

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"gtithub.com/jgfranco17/opsrunner/cli/config"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
)

var _ = Describe("User Stories Acceptance Tests", func() {
	var tempDir string
	var realExecutor *executor.DefaultExecutor
	var ctx context.Context

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		realExecutor = &executor.DefaultExecutor{}
		ctx = context.Background()
	})

	Describe("User Story: As a developer, I want to build my project", func() {
		Context("Given I have a valid project configuration", func() {
			It("When I run the build command, Then my project should build successfully", func() {
				// Given: A developer has a valid project configuration
				projectConfig := &config.ProjectDefinition{
					Name:        "DeveloperProject",
					Version:     "1.0.0",
					Description: "A developer's project",
					RepoUrl:     "https://github.com/developer/project",
					Codebase: config.Codebase{
						Language:     "go",
						Dependencies: "go.mod",
						Install: config.Operation{
							Steps: []string{
								"echo 'Installing dependencies...'",
								"go mod tidy",
							},
						},
						Build: config.Operation{
							Steps: []string{
								"echo 'Building project...'",
								"echo 'go build -o myapp'",
								"echo 'Build completed successfully'",
							},
						},
					},
				}

				// And: The configuration file exists
				content, err := yaml.Marshal(projectConfig)
				Expect(err).To(BeNil())
				configPath := filepath.Join(tempDir, ".opsrunner.yaml")
				err = os.WriteFile(configPath, content, 0644)
				Expect(err).To(BeNil())

				// When: The developer runs the build command
				buildOptions := &config.BuildOptions{
					NoInstall: false,
				}
				err = config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The project should build successfully
				Expect(err).To(BeNil())
			})
		})

		Context("Given I want to skip dependency installation", func() {
			It("When I run the build command with no-install flag, Then only build steps should execute", func() {
				// Given: A developer wants to skip dependency installation
				projectConfig := &config.ProjectDefinition{
					Name:    "SkipInstallProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Install: config.Operation{
							Steps: []string{
								"echo 'This should be skipped'",
								"go mod download",
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

				// When: The developer runs the build command with no-install flag
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: Only build steps should execute
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("User Story: As a developer, I want to handle build failures gracefully", func() {
		Context("Given a build step fails", func() {
			It("When fail_fast is enabled, Then execution should stop immediately", func() {
				// Given: A build step will fail
				projectConfig := &config.ProjectDefinition{
					Name:    "FailingBuildProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Build: config.Operation{
							FailFast: true,
							Steps: []string{
								"echo 'Step 1 - This will succeed'",
								"exit 1", // This will fail
								"echo 'Step 3 - This should not execute'",
							},
						},
					},
				}

				// When: The build is executed
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: Execution should stop at the failing step
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to run build steps"))
			})

			It("When fail_fast is disabled, Then all steps should execute and report failures", func() {
				// Given: A build step will fail but fail_fast is disabled
				projectConfig := &config.ProjectDefinition{
					Name:    "NonFailFastProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Build: config.Operation{
							FailFast: false,
							Steps: []string{
								"echo 'Step 1 - This will succeed'",
								"exit 1", // This will fail
								"echo 'Step 3 - This should still execute'",
							},
						},
					},
				}

				// When: The build is executed
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: All steps should execute but an error should be reported
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("failed to run steps"))
			})
		})
	})

	Describe("User Story: As a developer, I want to set environment variables for my build", func() {
		Context("Given I have environment variables in my configuration", func() {
			It("When I run the build, Then environment variables should be available", func() {
				// Given: A project configuration with environment variables
				projectConfig := &config.ProjectDefinition{
					Name:    "EnvVarProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Build: config.Operation{
							Env: map[string]string{
								"GO_ENV":      "production",
								"CGO_ENABLED": "0",
								"BUILD_TAG":   "v1.0.0",
							},
							Steps: []string{
								"echo 'Environment variables should be set'",
								"echo 'go build'",
							},
						},
					},
				}

				// When: The build is executed
				buildOptions := &config.BuildOptions{
					NoInstall: true,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully with environment variables
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("User Story: As a DevOps engineer, I want to integrate OpsRunner into CI/CD pipelines", func() {
		Context("Given I have a CI/CD pipeline configuration", func() {
			It("When I run OpsRunner in the pipeline, Then it should provide consistent builds", func() {
				// Given: A CI/CD pipeline configuration
				projectConfig := &config.ProjectDefinition{
					Name:    "CICDProject",
					Version: "1.0.0",
					Codebase: config.Codebase{
						Language: "go",
						Install: config.Operation{
							FailFast: true,
							Env: map[string]string{
								"CI": "true",
							},
							Steps: []string{
								"go mod download",
								"go mod verify",
							},
						},
						Build: config.Operation{
							FailFast: true,
							Env: map[string]string{
								"CI":          "true",
								"CGO_ENABLED": "0",
								"GOOS":        "linux",
								"GOARCH":      "amd64",
							},
							Steps: []string{
								"echo 'go test ./...'",
								"echo 'go build -o app'",
								"echo 'CI/CD build completed'",
							},
						},
					},
				}

				// When: OpsRunner is executed in the CI/CD pipeline
				buildOptions := &config.BuildOptions{
					NoInstall: false,
				}
				err := config.Build(ctx, realExecutor, projectConfig, buildOptions)

				// Then: The build should complete successfully for CI/CD
				Expect(err).To(BeNil())
			})
		})
	})
})
