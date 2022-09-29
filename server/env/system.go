package env

import (
	"log"
	"os"
	"path/filepath"
)

type system interface {
	SpecificSetup()
	GetInstallKeyPath() string
}

func GetHomeDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	if len(homePath) < 1 {
		log.Fatalln("The user home directory path is not valid: " + string(homePath))
	}
	return homePath
}

func GetContextAppPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/bin/context_app/")
}
