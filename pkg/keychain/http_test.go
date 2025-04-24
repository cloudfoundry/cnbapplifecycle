package keychain_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/keychain"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP RoundTripper", func() {
	var (
		creds  authn.Keychain
		client *http.Client
		err    error
	)

	BeforeEach(func() {
		creds, err = keychain.FromEnv()
		Expect(err).ToNot(HaveOccurred())
		client = keychain.NewHTTPClient(creds)
	})

	Describe("RoundTrip", func() {
		It("works when keychain is nil", func() {
			httpmock.RegisterResponder("GET", "https://test.io", func(r *http.Request) (*http.Response, error) {
				if a := r.Header.Get("Authorization"); a != "" {
					return httpmock.NewStringResponse(http.StatusBadRequest, fmt.Sprintf("found authorization header: %q", a)), nil
				}
				return httpmock.NewStringResponse(http.StatusOK, ""), nil
			})

			client = keychain.NewHTTPClient(nil)

			res, err := client.Get("https://test.io")
			Expect(err).ToNot(HaveOccurred())
			defer res.Body.Close()
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		It("works with the default keychain (anonymous)", func() {
			httpmock.RegisterResponder("GET", "https://test.io", func(r *http.Request) (*http.Response, error) {
				if a := r.Header.Get("Authorization"); a != "" {
					return httpmock.NewStringResponse(http.StatusBadRequest, fmt.Sprintf("found authorization header: %q", a)), nil
				}
				return httpmock.NewStringResponse(http.StatusOK, ""), nil
			})
			res, err := client.Get("https://test.io")
			Expect(err).ToNot(HaveOccurred())
			defer res.Body.Close()
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		Describe("with credentials", func() {
			BeforeEach(func() {
				Expect(os.Setenv(keychain.CnbCredentialsEnv, `{"bearer.test":{"token":"foo"},"basic.test":{"username":"foo","password":"bar"}}`)).To(Succeed())
				creds, err = keychain.FromEnv()
				Expect(err).ToNot(HaveOccurred())
				client = keychain.NewHTTPClient(creds)
			})

			AfterEach(func() {
				Expect(os.Unsetenv(keychain.CnbCredentialsEnv)).To(Succeed())
			})

			It("sets Authorization header (token)", func() {
				httpmock.RegisterResponder("GET", "https://bearer.test", func(r *http.Request) (*http.Response, error) {
					authHeader := r.Header.Get("Authorization")
					if authHeader == "" {
						return httpmock.NewStringResponse(http.StatusBadRequest, "authorization header not found"), nil
					}

					return httpmock.NewStringResponse(http.StatusOK, authHeader), nil
				})

				res, err := client.Get("https://bearer.test")
				Expect(err).ToNot(HaveOccurred())
				defer res.Body.Close()
				Expect(res.StatusCode).To(Equal(http.StatusOK))

				body := bytes.NewBuffer(nil)
				_, err = io.Copy(body, res.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(body.String()).To(Equal("Bearer foo"))
			})

			It("sets Authorization header (basic)", func() {
				httpmock.RegisterResponder("GET", "https://basic.test", func(r *http.Request) (*http.Response, error) {
					authHeader := r.Header.Get("Authorization")
					if authHeader == "" {
						return httpmock.NewStringResponse(http.StatusBadRequest, "authorization header not found"), nil
					}

					return httpmock.NewStringResponse(http.StatusOK, authHeader), nil
				})

				res, err := client.Get("https://basic.test")
				Expect(err).ToNot(HaveOccurred())
				defer res.Body.Close()
				Expect(res.StatusCode).To(Equal(http.StatusOK))

				body := bytes.NewBuffer(nil)
				_, err = io.Copy(body, res.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(body.String()).To(Equal("Basic Zm9vOmJhcg=="))
			})
		})
	})
})
