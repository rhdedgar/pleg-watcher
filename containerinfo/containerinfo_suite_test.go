package containerinfo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestContainerinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containerinfo Suite")
}
