package buildpacks_test

import (
	"fmt"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/buildpacks"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/log"

	"github.com/cespare/xxhash/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Translate", func() {
	var err error
	var bpDir string
	var bps []string
	var hashedName string
	var logger = log.NewLogger()

	BeforeEach(func() {
		bpDir, err = os.MkdirTemp("", "buildpacks")
		Expect(err).NotTo(HaveOccurred())

		bps = []string{"foo", "bar"}
		hashedName = fmt.Sprintf("%016x", xxhash.Sum64String("foo"))
		Expect(os.MkdirAll(filepath.Join(bpDir, hashedName), 0o755)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(bpDir)).To(Succeed())
	})

	It("correctly archives downloaded system-buildpacks and translates the uris", func() {
		bps, err = buildpacks.Translate(bps, bpDir, logger)
		Expect(err).NotTo(HaveOccurred())
		Expect(bps).To(Equal([]string{
			fmt.Sprintf("file://%s", filepath.Join(bpDir, hashedName)),
			"bar",
		}))
		Expect(filepath.Join(bpDir, hashedName)).To(BeADirectory())
	})

	When("when hashed path is not a directory", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll(filepath.Join(bpDir, hashedName))).To(Succeed())
			Expect(os.WriteFile(filepath.Join(bpDir, hashedName), []byte("some file"), 0o644)).To(Succeed())
		})

		It("throws an error", func() {
			bps, err = buildpacks.Translate(bps, bpDir, logger)
			Expect(bps).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("is not a directory")))
		})
	})
})
