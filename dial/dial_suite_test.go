package dial_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDial(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dial Suite")
}
