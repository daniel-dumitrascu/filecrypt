package env

import (
	"os"
	"server/utils"
)

type system interface {
	SpecificSetup()
	GetInterpretor() string
	GetBinDirPath() string
	GetKeysDirPath() string
	ChangeFilePermission(keyPath *string)
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
