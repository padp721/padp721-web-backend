package utilities

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func LoadPublicKey(privateKey *rsa.PrivateKey) crypto.PublicKey {
	return privateKey.Public()
}
