package archive_test

import (
	"archive/tar"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/archive"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type fakeWriter struct {
	headers []*tar.Header
}

func (fw *fakeWriter) WriteHeader(hdr *tar.Header) error {
	fw.headers = append(fw.headers, hdr)
	return nil
}

func (fw *fakeWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (fw *fakeWriter) Close() error {
	return nil
}

var _ = Describe("FromDirectory", func() {
	It("packages a droplet", func() {
		writer := &fakeWriter{}
		Expect(archive.FromDirectory("./testdata", writer)).To(Succeed())

		Expect(writer.headers).To(HaveLen(6))

		Expect(writer.headers[0].Name).To(Equal("bar"))
		Expect(writer.headers[0].Typeflag).To(Equal(uint8(tar.TypeReg)))

		Expect(writer.headers[1].Name).To(Equal("foo"))
		Expect(writer.headers[1].Typeflag).To(Equal(uint8(tar.TypeReg)))

		Expect(writer.headers[2].Name).To(Equal("foobar"))
		Expect(writer.headers[2].Typeflag).To(Equal(uint8(tar.TypeDir)))

		Expect(writer.headers[3].Name).To(Equal("foobar/bazz"))
		Expect(writer.headers[3].Typeflag).To(Equal(uint8(tar.TypeReg)))

		Expect(writer.headers[4].Name).To(Equal("link"))
		Expect(writer.headers[4].Typeflag).To(Equal(uint8(tar.TypeSymlink)))
		Expect(writer.headers[4].Linkname).To(Equal("foobar/bazz"))

		Expect(writer.headers[5].Name).To(Equal("templink"))
		Expect(writer.headers[5].Typeflag).To(Equal(uint8(tar.TypeSymlink)))
		Expect(writer.headers[5].Linkname).To(Equal("/tmp/test"))
	})
})
