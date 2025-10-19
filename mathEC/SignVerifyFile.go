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
	"strings"
)

type EcdsaSignature struct {
	R, S *big.Int
}

func SignDocument(filePath string, key *ecdsa.PrivateKey) (string, error) {
	msgToSign, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}
	// HASH(m) usando SHA-256 CheckSum
	hash := sha256.Sum256(msgToSign)
	// firma hash con privada, genera (r,s)
	r, s, err := ecdsa.Sign(rand.Reader, key, hash[:])
	if err != nil {
		return "", fmt.Errorf("error al firmar el archivo: %v", err)
	}

	// Marshalling de firma en formato ASN.1/DER para archivo
	signBytes, err := asn1.Marshal(EcdsaSignature{R: r, S: s})
	if err != nil {
		return "", fmt.Errorf("error al hacer marshal del archivo: %v", err)
	}
	newName := filePath + ".signed"
	file, err := os.Create(newName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo: %v", err)
	}
	// al salir de la func cierra archivo
	defer file.Close()
	// longitud firma
	sigLen := uint32(len(signBytes))
	// escribir longitud firma
	err = binary.Write(file, binary.BigEndian, sigLen)
	if err != nil {
		return "", fmt.Errorf("error al escribir longitud de firma: %v", err)
	}
	// escribir firma
	_, err = file.Write(signBytes)
	if err != nil {
		return "", fmt.Errorf("error al escribir firma: %v", err)
	}
	// añadir archivo original
	_, err = file.Write(msgToSign)
	if err != nil {
		return "", fmt.Errorf("error al escribir mensaje: %v", err)
	}
	// regresa nombre nuevo y no error
	return newName, nil
}

func VerifyFile(filePath string, key *ecdsa.PublicKey) error {
	combinedData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo original: %v", err)
	}
	var extractedSigLen uint32
	// leer longitud de firma en primeros 4 bytes (uint32 -> log2(32) = 4)
	reader := bytes.NewReader(combinedData[0:4])
	err = binary.Read(reader, binary.BigEndian, &extractedSigLen)
	if err != nil {
		return fmt.Errorf("error al leer la longitud de la firma: %v", err)
	}
	// extraer firma en los sig. n bytes
	sigStart := uint32(4)
	sigEnd := sigStart + extractedSigLen
	extractedSigBytes := combinedData[sigStart:sigEnd]

	// extraer lo que queda que es el mensaje
	extractedMessage := combinedData[sigEnd:]

	// unmarshal de firma
	var extractedSig EcdsaSignature
	_, err = asn1.Unmarshal(extractedSigBytes, &extractedSig)
	if err != nil {
		return fmt.Errorf("error al hacer unmarshal de firma: %v", err)
	}
	// HASH(m)
	extractedHash := sha256.Sum256(extractedMessage)
	// verificación
	isValid := ecdsa.Verify(key, extractedHash[:], extractedSig.R, extractedSig.S)
	if !isValid {
		return fmt.Errorf("la firma no es valida o ha sido alterada")
	}
	prevNameFile := strings.TrimSuffix(filePath, ".signed")
	_, err = os.Stat(prevNameFile)
	if os.IsNotExist(err) {
		err = os.WriteFile(prevNameFile, extractedMessage, 0644)
		if err != nil {
			return fmt.Errorf("error escribiendo el archivo original: %v", err)
		}
	}
	return nil
}
