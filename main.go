package main

import (
	filepick "ecdsa_gui/gui/file-pick"
	getstring "ecdsa_gui/gui/get-string"
	mainmenu "ecdsa_gui/gui/main-menu"
	showdialog "ecdsa_gui/gui/show-dialog"
	"ecdsa_gui/logging"
	mathec "ecdsa_gui/mathEC"
	"flag"
)

var (
	privKeyDefined = false
	pubKeyDefined  = false
	keys           mathec.Keys
)

func main() {
	// preparar para logging
	logPtr := flag.Bool("log", false, "produces logging output of the program")
	flag.Parse()
	if *logPtr {
		err := logging.StartLogging()
		if err != nil {
			showdialog.ShowError(err.Error())
			return
		}
		defer logging.LogService.Close()
	}
	for {
		// hasta que el usuario decida salirse, menú
		n := mainmenu.ObtenerOpcion()
		switch n {
		// salir
		case mainmenu.ExitSelected:
			return
		// generar llaves
		case mainmenu.GenerateKeysSelected:
			//INPUT: nombre con el que se identifican las llaves
			nameKeys, exit := getstring.GetNameGeneratedKeys()
			if exit {
				continue
			}
			//OUTPUT: llaves .pem (curva P-256)
			keyGen, err := mathec.GenerateKeys(nameKeys)
			if err != nil {
				showdialog.ShowError(err.Error())
				continue
			}
			showdialog.ShowDialog("¡Llaves guardadas con éxito!", 3)
			// guardamos las llaves para no cargarlas de nuevo en caso de
			// continuar en la misma sesión
			privKeyDefined = true
			pubKeyDefined = true
			keys = keyGen
		case mainmenu.SignSelected:
			// si no se han cargado las llaves, las cargamos.
			if !privKeyDefined {
				//FILEPICKER: llave privada
				privKeyPath, exit := filepick.GetKeyFile("Llave privada")
				if exit {
					continue
				}
				privKey, err := mathec.LoadPrivateKey(privKeyPath)
				if err != nil {
					showdialog.ShowError(err.Error())
					continue
				}
				keys.PrivKey = privKey
				privKeyDefined = true
			}
			//FILEPICKER: archivo a firmar
			filePath, quit := filepick.GetFile("a firmar")
			if quit {
				continue
			}
			//firma de documento con ecdsa (sha256)
			//OUTPUT: archivo firmado .signed
			firmFileName, err := mathec.SignDocument(filePath, keys.PrivKey)
			if err != nil {
				showdialog.ShowError(err.Error())
				continue
			}
			showdialog.ShowDialog("¡Documento firmado con éxito en"+firmFileName+"!", 3)
		// Verificación de firma
		case mainmenu.VerifySelected:
			// si no se han cargado las llaves, las cargamos.
			if !pubKeyDefined {
				pubKeyPath, exit := filepick.GetKeyFile("Llave pública")
				if exit {
					continue
				}
				//INPUT: Llave pública
				pubKey, err := mathec.LoadPublicKey(pubKeyPath)
				if err != nil {
					showdialog.ShowError(err.Error())
					continue
				}
				keys.PubKey = pubKey
				pubKeyDefined = true
			}
			//FILEPICKER: Archivo firmado .signed
			filePath, quit := filepick.GetSignFile("a verificar")
			if quit {
				continue
			}
			err := mathec.VerifyFile(filePath, keys.PubKey)
			// si falla, muestra error específico
			if err != nil {
				showdialog.ShowError("La verificación falló: " + err.Error())
				continue
			}
			showdialog.ShowDialog("¡La firma es valida!", 3)
		}
	}
}
