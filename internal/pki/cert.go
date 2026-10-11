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

func IssueCertificate(domain string) (string, string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", "", err
	}

	// Load Root CA
	caCertPath := filepath.Join(dir, "rootCA.pem")
	caKeyPath := filepath.Join(dir, "rootCA.key")

	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return "", "", err
	}
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return "", "", err
	}

	block, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", "", err
	}

	blockKey, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(blockKey.Bytes)
	if err != nil {
		return "", "", err
	}

	// Generate Leaf Private Key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   domain,
			Organization: []string{"Dev-PKI Local Cert"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().AddDate(1, 0, 0), // 1 year validity
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{domain, "*." + domain, "localhost"},
		BasicConstraintsValid: true,
	}

	// Sign Leaf Cert with Root CA
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priv.PublicKey, caKey)
	if err != nil {
		return "", "", err
	}

	// Output files to current working directory
	certFileName := domain + ".crt"
	keyFileName := domain + ".key"

	certOut, err := os.Create(certFileName)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.OpenFile(keyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certFileName, keyFileName, nil
}