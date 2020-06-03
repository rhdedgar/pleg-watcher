package containerscan_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestContainerscan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containerscan Suite")
}
