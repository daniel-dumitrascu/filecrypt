package keymgn

import (
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"server/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

func LoadKey(installDir *string) []byte {
	// Take the key that has the right name format --> install_key_[timestamp]
	// having the latest timestamp
	var keyMatches []string
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("install_key_*", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			keyMatches = append(keyMatches, path)
		}
		return nil
	}

	log := utils.GetLogger()

	if err := filepath.Walk(*installDir, walkFunc); err != nil {
		log.Error("Error when trying to load the key: %s", err)
		return nil
	}

	if len(keyMatches) == 0 {
		log.Error("There is no key found to load")
		return nil
	}

	sort.Slice(keyMatches, func(i int, j int) bool {
		strTimestamp_a := keyMatches[i][strings.LastIndex(keyMatches[i], "_")+1 : len(keyMatches[i])]
		strTimestamp_b := keyMatches[j][strings.LastIndex(keyMatches[j], "_")+1 : len(keyMatches[j])]

		a, _ := strconv.ParseInt(strTimestamp_a, 10, 64)
		b, _ := strconv.ParseInt(strTimestamp_b, 10, 64)

		return a > b
	})

	log.Info("Loaded key: " + keyMatches[0])
	key, err := os.ReadFile(keyMatches[0])
	if err != nil {
		log.Error("Issue reading the key")
		return nil
	}

	key, err = base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		log.Error("Issue decoding the key")
		return nil
	}

	return key
}

func InstallKey(inputKeyPath *string, outputKeyPath *string) string {
	log := utils.GetLogger()
	inputFile, err := os.Open(*inputKeyPath)
	if err != nil {
		log.Error("Cannot install key. Key path is not valid: %s", err)
		return ""
	}

	outputFile, err := os.Create(*outputKeyPath)
	if err != nil {
		inputFile.Close()
		log.Error("Couldn't open dest file: %s", err)
		return ""
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		log.Error("Writing to output file failed: %s", err)
		os.Remove(*outputKeyPath)
		return ""
	}

	log.Info("The new installed key has been loaded: " + *outputKeyPath)
	return *outputKeyPath
}

func GenerateKeyName() string {
	return "install_key_" + strconv.FormatInt(time.Now().UnixMicro(), 10)
}
