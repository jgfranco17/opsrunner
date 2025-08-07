package cli_commands_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCliCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CliCommands Suite")
}
