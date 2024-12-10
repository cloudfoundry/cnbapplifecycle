package cli_test

import (
	"code.cloudfoundry.org/cnbapplifecycle/cmd/launcher/cli"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Launch", func() {

	Context("a valid metadata.toml", func() {
		BeforeEach(func() {
			GinkgoT().Setenv("CNB_LAYERS_DIR", "../../../integration/testdata/validMetadata")
		})

		It("launches an app process", func() {
			launcher := cli.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "python app.py"}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedCmd).To(Equal("python app.py"))
			Expect(launcher.ExecutedDirect).To(Equal(false))
			Expect(launcher.ExecutedType).To(BeEmpty())
		})

		It("launches a task", func() {
			launcher := cli.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "--", "python --version"}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedCmd).To(Equal("python --version"))
			Expect(launcher.ExecutedDirect).To(Equal(false))
			Expect(launcher.ExecutedType).To(BeEmpty())
		})

		It("launches a sidecar process", func() {
			launcher := cli.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "./sidecar.sh"}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedCmd).To(Equal("./sidecar.sh"))
			Expect(launcher.ExecutedDirect).To(Equal(false))
			Expect(launcher.ExecutedType).To(Equal("sidecar"))
		})
	})
})
