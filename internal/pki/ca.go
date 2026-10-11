package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".dev-pki")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

func CAExists() bool {
	dir, err := GetConfigDir()
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "rootCA.pem")); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "rootCA.key")); os.IsNotExist(err) {
		return false
	}
	return true
}

func GenerateRootCA() (string, string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", "", err
	}

	certPath := filepath.Join(dir, "rootCA.pem")
	keyPath := filepath.Join(dir, "rootCA.key")

	if CAExists() {
		return certPath, keyPath, nil
	}

	// 1. Generate Private Key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	// 2. Setup Root CA Certificate Template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Dev-PKI Local Root CA",
			Organization: []string{"Dev-PKI Local Development"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years validity
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
	}

	// 3. Self-sign the Root CA certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	// 4. Save Certificate PEM
	certOut, err := os.Create(certPath)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// 5. Save Private Key PEM
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	return certPath, keyPath, nil
}