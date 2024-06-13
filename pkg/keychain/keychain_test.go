package keychain_test

import (
	"os"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/keychain"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromEnv", func() {
	var registry, bearerRegistry name.Registry
	var anonymousRegistry name.Registry

	BeforeEach(func() {
		Expect(os.Setenv(keychain.CnbCredentialsEnv, `{"registry.io":{"username":"username","password":"password"}, "bearer.io":{"token":"token"}}`)).To(Succeed())

		registry, _ = name.NewRegistry("registry.io")
		bearerRegistry, _ = name.NewRegistry("bearer.io")
		anonymousRegistry, _ = name.NewRegistry("no-creds.io")
	})

	AfterEach(func() {
		Expect(os.Unsetenv(keychain.CnbCredentialsEnv)).To(Succeed())
	})

	Describe("Resolve", func() {
		It("resolves credentials for the given registries", func() {
			creds, err := keychain.FromEnv()
			Expect(err).ToNot(HaveOccurred())

			authentication, err := creds.Resolve(registry)
			Expect(err).ToNot(HaveOccurred())
			Expect(authentication).ToNot(BeNil())
			authorization, err := authentication.Authorization()
			Expect(err).ToNot(HaveOccurred())
			Expect(authorization).ToNot(BeNil())
			Expect(authorization).To(BeEquivalentTo(&authn.AuthConfig{
				Username: "username",
				Password: "password",
			}))

			bearerAuth, err := creds.Resolve(bearerRegistry)
			Expect(err).ToNot(HaveOccurred())
			Expect(bearerAuth).ToNot(BeNil())
			bearerAuthz, err := bearerAuth.Authorization()
			Expect(err).ToNot(HaveOccurred())
			Expect(bearerAuthz).ToNot(BeNil())
			Expect(bearerAuthz).To(BeEquivalentTo(&authn.AuthConfig{
				RegistryToken: "token",
			}))

			anonAuth, err := creds.Resolve(anonymousRegistry)
			Expect(err).ToNot(HaveOccurred())
			Expect(anonAuth).ToNot(BeNil())
			anonAuthz, err := anonAuth.Authorization()
			Expect(err).ToNot(HaveOccurred())
			Expect(anonAuthz).ToNot(BeNil())
			Expect(anonAuthz).To(BeEquivalentTo(&authn.AuthConfig{}))
		})

		DescribeTable("credentials combinations (success)",
			func(envVal string) {
				Expect(os.Setenv(keychain.CnbCredentialsEnv, envVal)).To(Succeed())

				creds, err := keychain.FromEnv()
				Expect(err).ToNot(HaveOccurred())

				authentication, err := creds.Resolve(registry)
				Expect(err).ToNot(HaveOccurred())
				Expect(authentication).ToNot(BeNil())
			},
			Entry("username and password", `{"registry.io":{"username":"username", "password":"password"}}`),
			Entry("token", `{"registry.io":{"token":"token"}}`),
		)

		DescribeTable("credentials combinations (failure)",
			func(envVal string) {
				Expect(os.Setenv(keychain.CnbCredentialsEnv, envVal)).To(Succeed())

				creds, err := keychain.FromEnv()
				Expect(err).ToNot(HaveOccurred())

				authentication, err := creds.Resolve(registry)
				Expect(err).To(MatchError("invalid credential combination"))
				Expect(authentication).To(BeNil())
			},
			Entry("no fields", `{"registry.io":{}}`),
			Entry("all fields", `{"registry.io": {"username":"username", "password":"password", "token":"token"}}`),
			Entry("username and token", `{"registry.io": {"username":"username", "token":"token"}}`),
			Entry("password and token", `{"registry.io": {"password":"password", "token":"token"}}`),
			Entry("only username", `{"registry.io":{"username":"username"}}`),
			Entry("only password", `{"registry.io": {"password":"password"}}`),
		)
	})
})
