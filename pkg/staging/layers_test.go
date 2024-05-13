package staging_test

import (
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/log"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/staging"
	"github.com/buildpacks/lifecycle/buildpack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apexLog "github.com/apex/log"
)

type mockHandler struct {
	entries []*apexLog.Entry
}

func (m *mockHandler) HandleLog(e *apexLog.Entry) error {
	if strings.HasPrefix(e.Message, "removing layer") {
		m.entries = append(m.entries, e)
	}

	return nil
}

var _ = Describe("RemoveBuildOnlyLayers", func() {
	var layersDir string
	var buildpacks []buildpack.GroupElement
	var logger *log.Logger
	var logHandler *mockHandler

	BeforeEach(func() {
		layersDir = GinkgoT().TempDir()
		launchToml := []byte(`[types]
build = false
launch = true`)
		buildLaunchToml := []byte(`[types]
build = true
launch = true`)
		buildToml := []byte(`[types]
build = true
launch = false`)
		buildpacks = []buildpack.GroupElement{
			{
				ID:      "launch",
				Version: "1.0.0",
				API:     "0.8",
			},
			{
				ID:      "launch-build",
				Version: "1.0.0",
				API:     "0.8",
			},
		}

		logHandler = &mockHandler{}
		logger = log.NewLogger()
		logger.Handler = logHandler

		Expect(os.MkdirAll(filepath.Join(layersDir, "launch", "launch-layer"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(layersDir, "launch", "launch-layer.toml"), launchToml, 0o644)).To(Succeed())

		Expect(os.MkdirAll(filepath.Join(layersDir, "launch-build", "launch-layer"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(layersDir, "launch-build", "launch-layer.toml"), buildLaunchToml, 0o644)).To(Succeed())

		Expect(os.MkdirAll(filepath.Join(layersDir, "launch-build", "build-layer"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(layersDir, "launch-build", "build-layer.toml"), buildToml, 0o644)).To(Succeed())
	})

	It("removes all layers without launch = true", func() {
		Expect(staging.RemoveBuildOnlyLayers(layersDir, buildpacks, log.NewLogger())).To(Succeed())
		Expect(logHandler.entries).To(HaveLen(1))
	})
})
