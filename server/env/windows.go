//go:build windows

package env

import (
	"log"
	"os"
	"os/exec"
	"strings"
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

	// Create directory in /etc/context-app where the key is going to be stored
	if _, err := os.Stat(GetHomeDir() + "\\etc"); os.IsNotExist(err) {
		if createDirErr := os.Mkdir(GetHomeDir()+"\\etc", os.ModePerm); createDirErr != nil {
			log.Fatalln("Cannot create directory in \"etc\" in home", createDirErr)
		}
	}
	if _, err := os.Stat(GetInstallKeyPath()); os.IsNotExist(err) {
		if createDirErr := os.Mkdir(GetHomeDir()+"\\etc\\context-app", os.ModePerm); createDirErr != nil {
			log.Fatalln("Cannot create directory in \"\\etc\\context-app\" in home", createDirErr)
		}
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

func (sys *windows) GetInterpretor() string {
	//Find the path to the python exec
	cmd := exec.Command("where", "python")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Python not found on the system", err)
	}

	interpretor := strings.Replace(string(output), "\n", "", -1)
	interpretor = strings.Replace(interpretor, "\r", "", -1)
	return interpretor
}

func GetOsManager() system {
	return new(windows)
}
