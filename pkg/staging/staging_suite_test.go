package staging_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStaging(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Staging Suite")
}
