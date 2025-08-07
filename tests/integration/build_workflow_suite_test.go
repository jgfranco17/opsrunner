package build_workflow_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBuildWorkflow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BuildWorkflow Suite")
}
