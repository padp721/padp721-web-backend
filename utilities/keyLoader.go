package utilities

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func LoadPrivateKey(filename string) *rsa.PrivateKey {
	privateKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read Private Key: %v", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Failed to decode PEM block containing private key!")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return privateKey
}

func LoadPublicKey(privateKey *rsa.PrivateKey) crypto.PublicKey {
	return privateKey.Public()
}
