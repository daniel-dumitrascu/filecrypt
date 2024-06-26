//go:build windows

package env

import (
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"server/config"
	"server/utils"

	"golang.org/x/sys/windows/registry"
)

type windows struct {
}

func (sys *windows) SpecificSetup() {
	config.CURRENT_PLATFORM = config.PLATFORM_WIN

	keyEncryptName := "FilecryptEncrypt"
	keyDecryptName := "FilecryptDecrypt"
	keyAddKey := "FilecryptAddKey"
	keyGenKey := "FilecryptGenKey"
	execAppPath := sys.GetBinDirPath() + "\\" + config.APP_CLIENT_NAME + ".exe"
	encryptIconPath := sys.GetBinDirPath() + "\\..\\resources\\encrypt.ico"
	decryptIconPath := sys.GetBinDirPath() + "\\..\\resources\\decrypt.ico"
	keyIconPath := sys.GetBinDirPath() + "\\..\\resources\\key.ico"
	fileKeysPath := "*\\shell\\"
	genKeyPath := "Directory\\Background\\shell\\"
	dirKeysPath := "Folder\\shell\\"

	// Action keys for handling files
	if !IsKeyPresent(keyEncryptName, fileKeysPath) {
		CreateContextEntry(fileKeysPath, keyEncryptName, "Encrypt source", execAppPath, "encrypt", "%1", encryptIconPath)
	}

	if !IsKeyPresent(keyDecryptName, fileKeysPath) {
		CreateContextEntry(fileKeysPath, keyDecryptName, "Decrypt source", execAppPath, "decrypt", "%1", decryptIconPath)
	}

	if !IsKeyPresent(keyAddKey, fileKeysPath) {
		CreateContextEntry(fileKeysPath, keyAddKey, "Add key", execAppPath, "addkey", "%1", keyIconPath)
	}

	// Action keys for handling directories
	if !IsKeyPresent(keyEncryptName, dirKeysPath) {
		CreateContextEntry(dirKeysPath, keyEncryptName, "Encrypt source", execAppPath, "encrypt", "%1", encryptIconPath)
	}

	if !IsKeyPresent(keyDecryptName, dirKeysPath) {
		CreateContextEntry(dirKeysPath, keyDecryptName, "Decrypt source", execAppPath, "decrypt", "%1", decryptIconPath)
	}

	// Action for generating a new key
	if !IsKeyPresent(keyGenKey, genKeyPath) {
		CreateContextEntry(genKeyPath, keyGenKey, "Generate key", execAppPath, "genkey", "%v", keyIconPath)
	}
}

func IsKeyPresent(keyName string, path string) bool {
	log := utils.GetLogger()
	_, err := registry.OpenKey(registry.CLASSES_ROOT, path+keyName, registry.QUERY_VALUE)
	if err == syscall.ERROR_FILE_NOT_FOUND {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}

func CreateContextEntry(path string, contextName string, contextDesc string, appToExec string, action string, exeParam string, iconPath string) {
	log := utils.GetLogger()
	encryptKeyHandler, _, err := registry.CreateKey(registry.CLASSES_ROOT, path+contextName,
		registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		log.Fatal(err)
	}
	defer encryptKeyHandler.Close()

	err = encryptKeyHandler.SetStringValue("", contextDesc)
	if err != nil {
		log.Fatal(err)
	}

	err = encryptKeyHandler.SetStringValue("Icon", iconPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create command sub-key
	encryptSubKeyHandler, _, err := registry.CreateKey(registry.CLASSES_ROOT, path+contextName+"\\command",
		registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		log.Fatal(err)
	}
	defer encryptSubKeyHandler.Close()

	err = encryptSubKeyHandler.SetStringValue("", "\""+appToExec+"\""+" \""+action+"\" "+"\""+exeParam+"\"")
	if err != nil {
		log.Fatal(err)
	}
}

func (sys *windows) GetInterpretor() string {
	log := utils.GetLogger()
	//Find the path to the python exec
	cmd := exec.Command("where", "python")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Python not found on the system.", err)
	}

	lines := strings.Fields(string(output))
	return lines[0]
}

func (sys *windows) GetBinDirPath() string {
	return filepath.Join("C:/Program Files/" + config.APP_GENERIC_NAME + "/bin")
}

func (sys *windows) ChangeFilePermission(keyPath *string) {
	//no implementation needed
}

func (sys *windows) GetKeysDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/"+config.APP_GENERIC_NAME+"/keys")
}

func GetOsManager() system {
	return new(windows)
}
