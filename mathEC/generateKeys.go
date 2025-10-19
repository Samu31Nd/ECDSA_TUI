package mathec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

type Keys struct {
	PrivKey *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey
}

// generación de llaves usando curva P-256
func GenerateKeys(namePrefix string) (Keys, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return Keys{}, err
	}
	publicKey := &privateKey.PublicKey

	// estándar usando x509, en específico, PKCS#8
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return Keys{}, err
	}
	// estándar usando x509, en específico, PKIX
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return Keys{}, err
	}
	if err := writePEM("keys/"+namePrefix+"_priv.pem", "ECDSA PRIVATE KEY", privBytes, 0600); err != nil {
		return Keys{}, fmt.Errorf("writing private pem: %v", err)
	}

	if err := writePEM("keys/"+namePrefix+"_pub.pem", "ECDSA PUBLIC KEY", pubBytes, 0600); err != nil {
		return Keys{}, fmt.Errorf("writing private pem: %v", err)
	}

	return Keys{PubKey: publicKey, PrivKey: privateKey}, nil
}
