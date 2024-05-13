package staging_test

import (
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/staging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateEnvFiles", func() {
	var platformDir string

	BeforeEach(func() {
		GinkgoT().Setenv("FOO", "foo")
		GinkgoT().Setenv("BAR", "bar")
		GinkgoT().Setenv("BAZZ", "bazz")

		platformDir = GinkgoT().TempDir()
	})

	It("writes env files", func() {
		Expect(staging.CreateEnvFiles(platformDir, []string{"FOO", "BAR"})).To(Succeed())

		Expect(filepath.Join(platformDir, "env")).To(BeADirectory())

		entries, err := os.ReadDir(filepath.Join(platformDir, "env"))
		Expect(err).ToNot(HaveOccurred())
		Expect(entries).To(HaveLen(2))
		Expect(entries[0].Name()).To(Equal("BAR"))

		barContent, err := os.ReadFile(filepath.Join(platformDir, "env", "BAR"))
		Expect(err).ToNot(HaveOccurred())
		Expect(barContent).To(Equal([]byte("bar")))

		Expect(entries[1].Name()).To(Equal("FOO"))
		fooContent, err := os.ReadFile(filepath.Join(platformDir, "env", "FOO"))
		Expect(err).ToNot(HaveOccurred())
		Expect(fooContent).To(Equal([]byte("foo")))
	})

	It("throws an error if requested env var is not present", func() {
		Expect(staging.CreateEnvFiles(platformDir, []string{"NOT_PRESENT"})).To(MatchError(ContainSubstring("NOT_PRESENT")))
	})
})
