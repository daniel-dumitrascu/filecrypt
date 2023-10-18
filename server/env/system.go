package env

import (
	"log"
	"os"
	"path/filepath"
	"server/config"
)

type system interface {
	SpecificSetup()
	GetInterpretor() string
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

func GetAppDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/"+config.App_generic_name)
}

func GetBinDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/"+config.App_generic_name+"/"+config.App_bin_dir)
}

func GetKeysDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/"+config.App_generic_name+"/"+config.App_keys_dir)
}
