//go:build windows

package env

import (
	"log"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type windows struct {
}

func (sys *windows) SpecificSetup() {
	keyEncryptName := "FilecryptEncrypt"
	keyDecryptName := "FilecryptDecrypt"
	keyAddKey := "FilecryptAddKey"
	execAppPath := GetContextAppPath() + "\\context_app.exe"

	if !IsKeyPresent(keyEncryptName) {
		CreateContextEntry(keyEncryptName, "Encrypt source", execAppPath, "encrypt")
	}

	if !IsKeyPresent(keyDecryptName) {
		CreateContextEntry(keyDecryptName, "Decrypt source", execAppPath, "decrypt")
	}

	if !IsKeyPresent(keyAddKey) {
		CreateContextEntry(keyAddKey, "Add key", execAppPath, "addkey")
	}
}

func IsKeyPresent(keyName string) bool {
	keyPath := "*\\shell\\"
	_, err := registry.OpenKey(registry.CLASSES_ROOT, keyPath+keyName, registry.QUERY_VALUE)
	if err == syscall.ERROR_FILE_NOT_FOUND {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}

func CreateContextEntry(contextName string, contextDesc string, appToExec string, action string) {
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

	err = encryptSubKeyHandler.SetStringValue("", "\""+appToExec+"\""+" \""+action+"\" "+"\"%1\"")
	if err != nil {
		log.Fatal(err)
	}
}

func (sys *windows) GetInstallKeyPath() string {
	return GetHomeDir() + "\\etc\\context-app\\"
}

func GetOsManager() system {
	return new(windows)
}
