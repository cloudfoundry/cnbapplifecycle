package databaseuri_test

import (
	"testing"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/databaseuri"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDatabaseuri(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databaseuri Suite")
}

var _ = Describe("ParseDatabaseURI", func() {
	It("ignores services without credentials.uri", func() {
		services := `{"eg":[{}]}`
		uri, err := databaseuri.ParseDatabaseURI(services)
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(BeEmpty())
	})

	It("returns empty when there are non relational database services", func() {
		services := `{"eg":[{"credentials":{"uri":"sendgrid://foo:bar@host/db"}}]}`
		uri, err := databaseuri.ParseDatabaseURI(services)
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(BeEmpty())
	})

	It("returns empty when there are no services", func() {
		services := "{}"
		uri, err := databaseuri.ParseDatabaseURI(services)
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(BeEmpty())
	})

	It("returns empty when services are empty", func() {
		services := ""
		uri, err := databaseuri.ParseDatabaseURI(services)
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(BeEmpty())
	})

	Context("when there are relational database services", func() {
		Context("with a mysql URI", func() {
			It("changes the scheme to mysql2", func() {
				services := `{"eg":[{"credentials":{"uri":"mysql://username:password@host/db"}}]}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal("mysql2://username:password@host/db"))
			})
		})
		Context("with a mysql2 URI", func() {
			It("returns the URI unchanged", func() {
				services := `{"eg":[{"credentials":{"uri":"mysql2://username:password@host/db"}}]}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal("mysql2://username:password@host/db"))
			})
		})
		Context("with a postgres URI", func() {
			It("returns the URI unchanged", func() {
				services := `{"eg":[{"credentials":{"uri":"postgres://username:password@host/db"}}]}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal("postgres://username:password@host/db"))
			})
		})
		Context("with a postgresql URI", func() {
			It("changes the scheme to postgres", func() {
				services := `{"eg":[{"credentials":{"uri":"postgresql://username:password@host/db"}}]}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal("postgres://username:password@host/db"))
			})
		})
		Context("with multiple relational database URIs", func() {
			It("returns the first one found", func() {
				services := `{
					"abc":[{"credentials":{"uri":"postgres://username:password@host/db1"}},
					{"credentials":{"uri":"postgres://username:password@host/db2"}}]
				}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal("postgres://username:password@host/db1"))
			})
		})
		Context("with an invalid URI", func() {
			It("returns an empty string", func() {
				services := `{"eg":[{"credentials":{"uri":"postgresql://invalid:password@host/%a"}}]}`
				uri, err := databaseuri.ParseDatabaseURI(services)
				Expect(err).NotTo(HaveOccurred())
				Expect(uri).To(Equal(""))
			})
		})
	})

	It("handles multiple services correctly", func() {
		services := `{
			"abc":[{"credentials":{"uri":"u1"}}],
			"def":[{"other":"data"}],
			"ghi":[{"credentials":{"other":"data"}}],
			"jkl":[{},{"credentials":{"uri":"mysql://username:password@host/db"}}]
		}`
		uri, err := databaseuri.ParseDatabaseURI(services)
		Expect(err).NotTo(HaveOccurred())
		Expect(uri).To(Equal("mysql2://username:password@host/db"))
	})
})
