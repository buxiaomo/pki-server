package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func (service *service) Signature(CN string, NotBefore, NotAfter time.Time) (crt, crtsha256, key, keysha256 string) {
	pk, _ := generatePublicKey()
	tmp := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keysha256 = SHA256(tmp)
	key = string(tmp)

	k8sCrt, _ := generateCert(CN, NotBefore, NotAfter, pk)
	tmp = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: k8sCrt})
	crt = string(tmp)
	crtsha256 = SHA256(tmp)
	return
}

func generatePublicKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func generateCert(CN string, NotBefore time.Time, NotAfter time.Time, publicKey *rsa.PrivateKey) ([]byte, error) {
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: CN,
		},
		NotBefore:             NotBefore,
		NotAfter:              NotAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	return x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &publicKey.PublicKey, publicKey)
}
