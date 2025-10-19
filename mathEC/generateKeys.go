package mathec

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"ecdsa_gui/logging"
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
		logging.SendLog("[Error] GenerateKeys: error al intentar generar las llaves: %v", err)
		return Keys{}, err
	}
	publicKey := &privateKey.PublicKey
	logging.SendLog("[Log] GenerateKeys: Llave generada \n\tName: %v\n\tD (valor privado):%v\n\tB:%v\n\tx,y:(%v,%v)\n\tn: %v\n\tp:%v).", privateKey.Params().Name, privateKey.D, privateKey.Params().B, privateKey.X, privateKey.Y, privateKey.Params().N, privateKey.Params().P)

	// estándar usando x509, en específico, PKCS#8
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		logging.SendLog("[Error] GenerateKeys: error al intentar hacer marshal de la privada: %v", err)
		return Keys{}, err
	}
	// estándar usando x509, en específico, PKIX
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logging.SendLog("[Error] GenerateKeys: error al intentar hacer marshal de la publica: %v", err)
		return Keys{}, err
	}
	if err := writePEM("keys/"+namePrefix+"_priv.pem", "ECDSA PRIVATE KEY", privBytes, 0600); err != nil {
		logging.SendLog("[Error] GenerateKeys: error al intentar guardar la llave privada: %v", err)
		return Keys{}, fmt.Errorf("writing private pem: %v", err)
	}

	if err := writePEM("keys/"+namePrefix+"_pub.pem", "ECDSA PUBLIC KEY", pubBytes, 0600); err != nil {
		logging.SendLog("[Error] GenerateKeys: error al intentar guardar la llave publica: %v", err)
		return Keys{}, fmt.Errorf("writing private pem: %v", err)
	}

	return Keys{PubKey: publicKey, PrivKey: privateKey}, nil
}
