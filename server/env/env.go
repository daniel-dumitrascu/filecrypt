package env

import (
	"errors"
	"log"
	"os"
	"server_app/keymgn"
)

func Setup() {
	osmanager := GetOsManager()

	if _, err := os.Stat(osmanager.GetContextAppPath()); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Context app is not installed.")
	}

	osmanager.SpecificSetup()
	var installKeyPath = osmanager.GetInstallKeyPath()
	keymgn.LoadKey(&installKeyPath)
}

func GetInstallKeyPath() string {
	osmanager := GetOsManager()
	return osmanager.GetInstallKeyPath()
}
