package user_stories_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserStories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserStories Suite")
}
