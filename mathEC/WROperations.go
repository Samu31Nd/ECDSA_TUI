package mathec

import (
	"crypto/ecdsa"
	"crypto/x509"
	"ecdsa_gui/logging"
	"encoding/pem"
	"fmt"
	"os"
)

const PathKeys = "keys"

func writePEM(filename, blockType string, data []byte, perm os.FileMode) error {
	if _, err := os.Stat(PathKeys); os.IsNotExist(err) {
		logging.SendLog("[Error] writePEM: error al intentar abrir el directorio : %v\n\tIntentando crear directorio...", err)
		err2 := os.Mkdir(PathKeys, 0755)
		if err2 != nil {
			logging.SendLog("[Error] writePEM: error al crear directorio: %v", err2)
			return err2
		}
	}
	block := &pem.Block{
		Type:  blockType,
		Bytes: data,
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, block)
}

func readPEM(filename string) (*pem.Block, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", filename)
	}
	return block, nil
}

func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	block, err := readPEM(filename)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave que seleccionaste no es una llave privada del tipo ECDSA, es del tipo %T", key)
	}
	return privKey, nil
}

func LoadPublicKey(filename string) (*ecdsa.PublicKey, error) {
	block, err := readPEM(filename)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("la llave que seleccionaste no es una llave pública del tipo ECDSA, es del tipo %T", key)
	}
	return pubKey, nil
}
