package credhub_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"testing"
)

var certDir string

func TestPlatform(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Credhub Suite")
}

var _ = BeforeSuite(func() {
	var err error
	certDir, err = setupCerts()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(certDir)
	Expect(err).NotTo(HaveOccurred())
})

func setupCerts() (string, error) {
	dir, err := os.MkdirTemp("", "test-certs")
	if err != nil {
		return "", err
	}

	err = os.Mkdir(filepath.Join(dir, "cacerts"), 0777)
	if err != nil {
		return "", err
	}

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"some-org"}, CommonName: "ca"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		IsCA:                  true,
	}
	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, caKey.Public(), caKey)
	if err != nil {
		return "", err
	}

	err = writeAsPem(filepath.Join(dir, "cacerts"), "ca", caCert, nil)
	if err != nil {
		return "", err
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	serverTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{Organization: []string{"some-org"}, CommonName: "server"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	serverCert, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, serverKey.Public(), caKey)
	if err != nil {
		return "", err
	}

	err = writeAsPem(dir, "server", serverCert, x509.MarshalPKCS1PrivateKey(serverKey))
	if err != nil {
		return "", err
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	clientTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{Organization: []string{"some-org"}, CommonName: "client"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	clientCert, err := x509.CreateCertificate(rand.Reader, &clientTemplate, &caTemplate, clientKey.Public(), caKey)
	if err != nil {
		return "", err
	}

	err = writeAsPem(dir, "client", clientCert, x509.MarshalPKCS1PrivateKey(clientKey))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func writeAsPem(dir string, name string, cert []byte, key []byte) error {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.crt", name)), certPEM, 0644)
	if err != nil {
		return err
	}

	if len(key) > 0 {
		keyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: key,
		})
		err = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.key", name)), keyPEM, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
