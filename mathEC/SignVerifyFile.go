package mathec

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"ecdsa_gui/logging"
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
		logging.SendLog("[Error] SignDocument: error al leer el archivo: %v", err)
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}
	// HASH(m) usando SHA-256 CheckSum
	hash := sha256.Sum256(msgToSign)
	logging.SendLog("[Log] SignDocument: Llave privada \n\tName: %v\n\tD (valor privado):%v\n\tB:%v\n\tx,y:(%v,%v)\n\tn: %v\n\tp:%v).", key.Params().Name, key.D, key.Params().B, key.X, key.Y, key.Params().N, key.Params().P)
	logging.SendLog("[Log] SignDocument: Hash del archivo (SHA-256): (%v)", hash)
	// firma hash con privada, genera (r,s)
	r, s, err := ecdsa.Sign(rand.Reader, key, hash[:])
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al firmar el archivo: %v", err)
		return "", fmt.Errorf("error al firmar el archivo: %v", err)
	}
	logging.SendLog("[Log] SignDocument: Parametros de firma ECDSA:\n\tR:(%v)\n\tS:(%v)", r, s)

	// Marshalling de firma en formato ASN.1/DER para archivo
	signBytes, err := asn1.Marshal(EcdsaSignature{R: r, S: s})
	logging.SendLog("[Log] SignDocument: Longitud de firma: (%v)\n\tFirma en bytes: (%v)", len(signBytes), signBytes)
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al intentar marshal del archivo: %v", err)
		return "", fmt.Errorf("error al hacer marshal del archivo: %v", err)
	}
	newName := filePath + ".signed"
	logging.SendLog("[Log] SignDocument: Nombre del archivo donde se guardara la firma: (%v)", newName)
	file, err := os.Create(newName)
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al crear el archivo de firma: %v", err)
		return "", fmt.Errorf("error al crear el archivo: %v", err)
	}
	// al salir de la func cierra archivo
	defer file.Close()
	// longitud firma
	sigLen := uint32(len(signBytes))
	// escribir longitud firma
	err = binary.Write(file, binary.BigEndian, sigLen)
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al escribir la longitud de la firma: %v", err)
		return "", fmt.Errorf("error al escribir longitud de firma: %v", err)
	}
	// escribir firma
	_, err = file.Write(signBytes)
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al escribir la firma: %v", err)
		return "", fmt.Errorf("error al escribir firma: %v", err)
	}
	// añadir archivo original
	_, err = file.Write(msgToSign)
	if err != nil {
		logging.SendLog("[Error] SignDocument: error al escribir el mensaje en el archivo: %v", err)
		return "", fmt.Errorf("error al escribir mensaje: %v", err)
	}
	// regresa nombre nuevo y no error
	logging.SendLog("[Log] SignDocument: Se completo el guardado de firma y archivo correctamente.")
	return newName, nil
}

func VerifyFile(filePath string, key *ecdsa.PublicKey) error {
	combinedData, err := os.ReadFile(filePath)
	if err != nil {
		logging.SendLog("[Error] VerifyFile: error al leer el archivo original: %v", err)
		return fmt.Errorf("error al leer el archivo original: %v", err)
	}
	var extractedSigLen uint32
	// leer longitud de firma en primeros 4 bytes (uint32 -> log2(32) = 4)
	reader := bytes.NewReader(combinedData[0:4])
	err = binary.Read(reader, binary.BigEndian, &extractedSigLen)
	if err != nil {
		logging.SendLog("[Error] VerifyFile: error al leer la longitud de la firma: %v", err)
		return fmt.Errorf("error al leer la longitud de la firma: %v", err)
	}
	logging.SendLog("[Log] SignDocument: Longitud de firma: (%v)", extractedSigLen)
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
		logging.SendLog("[Error] VerifyFile: error al hacer unmarshal de firma: %v", err)
		return fmt.Errorf("error al hacer unmarshal de firma: %v", err)
	}
	logging.SendLog("[Log] SignDocument: Firma ECDSA: \n\tR:(%v)\n\tS:(%v)", extractedSig.R, extractedSig.S)
	// HASH(m)
	extractedHash := sha256.Sum256(extractedMessage)
	logging.SendLog("[Log] SignDocument: Hash extraido del archivo original: (%v)", extractedHash)
	// verificación
	isValid := ecdsa.Verify(key, extractedHash[:], extractedSig.R, extractedSig.S)
	if !isValid {
		// TODO: Proveer mejor log aqui en especifico
		logging.SendLog("[Error] VerifyFile: el algoritmo ecdsa.Verify ha determinado que la firma no es valida con el hash que se obtuvo\n\tHa sido alterado el archivo de origen o la firma.")
		return fmt.Errorf("la firma no es valida o ha sido alterada")
	}
	prevNameFile := strings.TrimSuffix(filePath, ".signed")
	_, err = os.Stat(prevNameFile)
	if os.IsNotExist(err) {
		logging.SendLog("[Error] VerifyFile: El archivo no existe: %v\n\tIntentando crear archivo...", err)
		err = os.WriteFile(prevNameFile, extractedMessage, 0644)
		if err != nil {
			logging.SendLog("[Error] VerifyFile: error intentando escribir el archivo original: %v", err)
			return fmt.Errorf("error escribiendo el archivo original: %v", err)
		}
	}
	logging.SendLog("[Log] SignDocument: La firma ha sido verificada y coincide, el archivo es legitimo.")
	return nil
}
