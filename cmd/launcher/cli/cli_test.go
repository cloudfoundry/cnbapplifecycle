package cli_test

import (
	"code.cloudfoundry.org/cnbapplifecycle/cmd/launcher/cli"
	"code.cloudfoundry.org/cnbapplifecycle/cmd/launcher/cli/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Launch", func() {

	Context("a valid metadata.toml", func() {
		BeforeEach(func() {
			GinkgoT().Setenv("CNB_LAYERS_DIR", "../../../integration/testdata/validMetadata")
		})

		It("launches an app process", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "python app.py", ""}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedSelf).To(Equal("web"))
			Expect(launcher.ExecutedCmd).To(BeEmpty())
		})

		It("launches a second app process", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "PORT=9876 python app.py", "--verbose", ""}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedSelf).To(Equal("web2"))
			Expect(launcher.ExecutedCmd).To(BeEmpty())
		})

		It("launches a task", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "--", "python --version"}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedCmd).To(Equal("python --version"))
			Expect(launcher.ExecutedSelf).To(Equal("launcher"))
		})

		It("launches a sidecar process", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "./sidecar.sh", ""}, &launcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(launcher.ExecutedCmd).To(Equal("./sidecar.sh"))
			Expect(launcher.ExecutedSelf).To(Equal("launcher"))
		})

		It("fails when no process is provided", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "", ""}, &launcher)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("launching failed"))
		})
	})
	Context("an invalid metadata.toml", func() {
		BeforeEach(func() {
			GinkgoT().Setenv("CNB_LAYERS_DIR", "../../../integration/testdata/invalidMetadata")
		})

		It("launches an app process", func() {
			launcher := fake.FakeLifecycleLauncher{}
			err := cli.Launch([]string{"/tmp/launcher", "app", "python app.py", ""}, &launcher)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("launching failed"))
		})
	})
})
