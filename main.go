package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gtithub.com/jgfranco17/opsrunner/cli/core"
	"gtithub.com/jgfranco17/opsrunner/cli/executor"
)

const (
	projectName        = "opsrunner"
	projectDescription = "OpsRunner: Running local devops with ease."
)

var (
	version string = "0.0.0-dev.1"
)

func main() {
	executor := &executor.DefaultExecutor{}
	command := core.NewCommandRegistry(projectName, projectDescription, version)
	commandsList := []*cobra.Command{
		core.GetBuildCommand(executor),
	}
	command.RegisterCommands(commandsList)

	err := command.Execute()
	if err != nil {
		log.Error(err.Error())
	}
}
