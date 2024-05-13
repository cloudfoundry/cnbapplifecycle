package integration_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("INCLUDE_INTEGRATION_TESTS") != "true" {
		return
	}

	format.MaxLength = 0

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite", Label("integration"))
}
