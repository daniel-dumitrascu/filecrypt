//go:build windows

package env

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

type windows struct {
}

func (sys *windows) GetInstallKeyPath() string {
	//TODO
	return ""
}

func (sys *windows) SpecificSetup() {
	//TODO
	/*k, err := registry.OpenKey(registry.CLASSES_ROOT, "*\\shell\\blabla\\command", registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	v, _, _ := k.GetStringValue("")
	fmt.Println("Value from registry: " + v)*/

	CreateContextEntry("FilecryptEncrypt", "Encrypt source", "C:\\Program Files (x86)\\Notepad++\\notepad++.exe")
	CreateContextEntry("FilecryptDecrypt", "Decrypt source", "C:\\Program Files (x86)\\Notepad++\\notepad++.exe")
}

func CreateContextEntry(contextName string, contextDesc string, appToExec string) {
	keyPath := "*\\shell\\"

	encryptKeyHandler, _, err := registry.CreateKey(registry.CLASSES_ROOT, keyPath+"\\"+contextName,
		registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		log.Fatal(err)
	}
	defer encryptKeyHandler.Close()

	err = encryptKeyHandler.SetStringValue("", contextDesc)
	if err != nil {
		log.Fatal(err)
	}

	// Create command sub-key
	encryptSubKeyHandler, _, err := registry.CreateKey(registry.CLASSES_ROOT, keyPath+"\\"+contextName+"\\command",
		registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		log.Fatal(err)
	}
	defer encryptSubKeyHandler.Close()

	err = encryptSubKeyHandler.SetStringValue("", "\""+appToExec+"\" \"%1\"")
	if err != nil {
		log.Fatal(err)
	}
}

func GetOsManager() system {
	return new(windows)
}
