package staging_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/log"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/staging"
	"github.com/BurntSushi/toml"
	"github.com/buildpacks/lifecycle/api"
	"github.com/buildpacks/pack/pkg/archive"
	"github.com/buildpacks/pack/pkg/blob"
	"github.com/buildpacks/pack/pkg/dist"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type fakeBlob struct {
	id string
}

func (f *fakeBlob) Open() (io.ReadCloser, error) {
	buf := &bytes.Buffer{}
	descriptor := dist.BuildpackDescriptor{
		WithAPI:    api.MustParse("0.3"),
		WithInfo:   dist.ModuleInfo{ID: f.id, Version: "1.1.0"},
		WithStacks: []dist.Stack{{ID: "some.stack.id"}},
	}
	var err error
	if err = toml.NewEncoder(buf).Encode(descriptor); err != nil {
		return nil, err
	}

	tarBuilder := archive.TarBuilder{}

	tarBuilder.AddFile("buildpack.toml", 0644, time.Now(), buf.Bytes())
	tarBuilder.AddDir("bin", 0644, time.Now())
	tarBuilder.AddFile("bin/build", 0644, time.Now(), []byte("build-contents"))
	tarBuilder.AddFile("bin/detect", 0644, time.Now(), []byte("detect-contents"))

	return tarBuilder.Reader(archive.DefaultTarWriterFactory()), err
}

type fakeDownloader struct {
}

func (f fakeDownloader) Download(ctx context.Context, pathOrURI string) (blob.Blob, error) {
	id := strings.ReplaceAll(pathOrURI, "file:/", "")
	return &fakeBlob{id: id}, nil
}

var _ = Describe("DownloadBuildpacks", func() {
	var err error
	var logger *log.Logger = log.NewLogger()
	var downloader = fakeDownloader{}
	var orderFile *os.File
	var buildpacksDir string

	BeforeEach(func() {
		orderFile, err = os.CreateTemp("", "orderToml")
		Expect(err).NotTo(HaveOccurred())
		buildpacksDir, err = os.MkdirTemp("", "buildpackDir")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(orderFile.Name())).To(Succeed())
		Expect(os.RemoveAll(buildpacksDir)).To(Succeed())
	})

	It("creates empty order.toml for empty buildpack list", func() {
		err = staging.DownloadBuildpacks([]string{}, buildpacksDir, nil, downloader, orderFile, logger)

		Expect(err).ToNot(HaveOccurred())

		orderToml := staging.OrderTOML{}

		b, err := os.ReadFile(orderFile.Name())
		Expect(err).NotTo(HaveOccurred())

		_, err = toml.Decode(string(b), &orderToml)
		Expect(err).NotTo(HaveOccurred())

		Expect(orderToml.Order).To(HaveLen(1))
		Expect(orderToml.Order[0].Group).To(HaveLen(0))
	})

	It("creates order.toml and downloads a buildpacks", func() {
		err = staging.DownloadBuildpacks([]string{"file:/buildpack1", "file:/buildpack2"}, buildpacksDir, nil, downloader, orderFile, logger)

		Expect(err).ToNot(HaveOccurred())

		orderToml := staging.OrderTOML{}

		b, err := os.ReadFile(orderFile.Name())
		Expect(err).NotTo(HaveOccurred())

		_, err = toml.Decode(string(b), &orderToml)
		Expect(err).NotTo(HaveOccurred())

		Expect(orderToml.Order).To(HaveLen(1))
		Expect(orderToml.Order[0].Group).To(HaveLen(2))
		Expect(orderToml.Order[0].Group[0].ID).To(Equal("buildpack1"))
		Expect(orderToml.Order[0].Group[1].ID).To(Equal("buildpack2"))

		_, err = os.Stat(filepath.Join(buildpacksDir, "buildpack1"))
		Expect(err).NotTo(HaveOccurred())
		_, err = os.Stat(filepath.Join(buildpacksDir, "buildpack2"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("works for duplicated buildpacks", func() {
		err = staging.DownloadBuildpacks([]string{"file:/buildpack", "file:/buildpack"}, buildpacksDir, nil, downloader, orderFile, logger)

		Expect(err).ToNot(HaveOccurred())

		_, err = os.Stat(filepath.Join(buildpacksDir, "buildpack"))
		Expect(err).NotTo(HaveOccurred())
	})

})
