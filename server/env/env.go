package env

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func Setup() {
	osmanager := GetOsManager()

	if _, err := os.Stat(osmanager.GetContextAppPath()); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Context app is not installed.")
	}

	osmanager.SpecificSetup()
}

func InstallKey(keyPath *string) {
	inputFile, err := os.Open(*keyPath)
	if err != nil {
		fmt.Printf("Cannot install key. Key path is not valid: %s", err)
		return
	}

	osmanager := GetOsManager()
	newKeyFilenamePath := osmanager.GetInstallKeyPath() + "install_key_" + strconv.FormatInt(time.Now().UnixMicro(), 10)
	outputFile, err := os.Create(newKeyFilenamePath)
	if err != nil {
		inputFile.Close()
		fmt.Printf("Couldn't open dest file: %s", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		fmt.Printf("Writing to output file failed: %s", err)
		os.Remove(newKeyFilenamePath)
		return
	}
}
