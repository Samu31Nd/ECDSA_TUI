package mathec

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
)

type EcdsaSignature struct {
	R, S *big.Int
}

func SignDocument(filePath string, key *ecdsa.PrivateKey) (string, error) {
	msgToSign, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}

	hash := sha256.Sum256(msgToSign)
	r, s, err := ecdsa.Sign(rand.Reader, key, hash[:])
	if err != nil {
		return "", fmt.Errorf("error al firmar el archivo: %v", err)
	}
	signBytes, err := asn1.Marshal(EcdsaSignature{R: r, S: s})
	if err != nil {
		return "", fmt.Errorf("error al hacer marshal del archivo: %v", err)
	}
	newName := filePath + ".signed"
	file, err := os.Create(newName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()
	sigLen := uint32(len(signBytes))
	err = binary.Write(file, binary.BigEndian, sigLen)
	if err != nil {
		return "", fmt.Errorf("error al escribir longitud de firma: %v", err)
	}
	_, err = file.Write(signBytes)
	if err != nil {
		return "", fmt.Errorf("error al escribir firma: %v", err)
	}
	_, err = file.Write(msgToSign)
	if err != nil {
		return "", fmt.Errorf("error al escribir mensaje: %v", err)
	}
	return newName, nil
}

func VerifyFile(filePath string, key *ecdsa.PublicKey) error {
	combinedData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo original: %v", err)
	}
	var extractedSigLen uint32
	reader := bytes.NewReader(combinedData[0:4])
	err = binary.Read(reader, binary.BigEndian, &extractedSigLen)
	if err != nil {
		return fmt.Errorf("error al leer la longitud de la firma: %v", err)
	}
	sigStart := uint32(4)
	sigEnd := sigStart + extractedSigLen
	extractedSigBytes := combinedData[sigStart:sigEnd]

	extractedMessage := combinedData[sigEnd:]

	var extractedSig EcdsaSignature
	_, err = asn1.Unmarshal(extractedSigBytes, &extractedSig)
	if err != nil {
		return fmt.Errorf("error al hacer unmarshal de firma: %v", err)
	}
	//HASH(m)
	extractedHash := sha256.Sum256(extractedMessage)
	//verify
	isValid := ecdsa.Verify(key, extractedHash[:], extractedSig.R, extractedSig.S)
	if !isValid {
		return fmt.Errorf("la firma no es valida o ha sido alterada")
	}
	return nil
}
