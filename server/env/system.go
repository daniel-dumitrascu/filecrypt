package env

import (
	"os"
	"path/filepath"
	"server/config"
	"server/utils"
)

type system interface {
	SpecificSetup()
	GetInterpretor() string
	GetBinDirPath() string
}

func GetHomeDir() string {
	log := utils.GetLogger()
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if len(homePath) < 1 {
		log.Fatal("The user home directory path is not valid: " + string(homePath))
	}
	return homePath
}

func GetKeysDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/."+config.App_generic_name)
}
