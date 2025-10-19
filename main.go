package main

import (
	filepick "ecdsa_gui/gui/file-pick"
	getstring "ecdsa_gui/gui/get-string"
	mainmenu "ecdsa_gui/gui/main-menu"
	showdialog "ecdsa_gui/gui/show-dialog"
	mathec "ecdsa_gui/mathEC"
)

var (
	privKeyDefined = false
	pubKeyDefined  = false
	keys           mathec.Keys
)

func main() {
	for {
		n := mainmenu.ObtenerOpcion()
		switch n {
		case mainmenu.ExitSelected:
			return
		case mainmenu.GenerateKeysSelected:
			nameKeys, exit := getstring.GetNameGeneratedKeys()
			if exit {
				continue
			}
			keyGen, err := mathec.GenerateKeys(nameKeys)
			if err != nil {
				showdialog.ShowError(err.Error())
				continue
			}
			showdialog.ShowDialog("¡Llaves guardadas con éxito!", 3)
			privKeyDefined = true
			pubKeyDefined = true
			keys = keyGen
		case mainmenu.SignSelected:
			// si no se han cargado las llaves, las cargamos.
			if !privKeyDefined {
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
			filePath, quit := filepick.GetFile("a firmar")
			if quit {
				continue
			}
			firmFileName, err := mathec.SignDocument(filePath, keys.PrivKey)
			if err != nil {
				showdialog.ShowError(err.Error())
				continue
			}
			showdialog.ShowDialog("¡Documento firmado con éxito en"+firmFileName+"!", 3)
		case mainmenu.VerifySelected:
			// si no se han cargado las llaves, las cargamos.
			if !pubKeyDefined {
				pubKeyPath, exit := filepick.GetKeyFile("Llave pública")
				if exit {
					continue
				}
				pubKey, err := mathec.LoadPublicKey(pubKeyPath)
				if err != nil {
					showdialog.ShowError(err.Error())
					continue
				}
				keys.PubKey = pubKey
				pubKeyDefined = true
			}
			filePath, quit := filepick.GetSignFile("a verificar")
			if quit {
				continue
			}
			err := mathec.VerifyFile(filePath, keys.PubKey)
			if err != nil {
				showdialog.ShowError("La verificación falló: " + err.Error())
				continue
			}
			showdialog.ShowDialog("¡La firma es valida!", 3)
		}
	}
}
