package env

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/keymgn"
	"strings"
)

type EnvData struct {
	loadedKey   string
	interpretor string
}

func Setup(data *EnvData) {
	fmt.Printf("Setup the environment\n")
	if err := SetupAppDirs(); errors.Is(err, os.ErrNotExist) {
		os.Exit(1)
	}

	osmanager := GetOsManager()
	osmanager.SpecificSetup()
	var installKeyPath = GetKeysDirPath()
	data.loadedKey = keymgn.LoadKey(&installKeyPath)
	data.interpretor = osmanager.GetInterpretor()

	var handleEncryptAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		inputPath := getStringFromReqBody(req)
		inputPath = strings.Replace(inputPath, "\"", "", -1)
		fmt.Println("Encrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			fmt.Println(computeErr)
			return
		}

		if len(data.loadedKey) > 0 {
			fmt.Println("Loaded key: " + data.loadedKey)
			CallScript(&data.interpretor, &inputPath, &outputPath, &data.loadedKey, "encrypt")
		} else {
			fmt.Println("Cannot encrypt because no key has been found")
		}
	}

	var handleDecryptAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		var inputPath string = getStringFromReqBody(req)
		inputPath = strings.Replace(inputPath, "\"", "", -1)
		fmt.Println("Decrypt action was triggered for: " + getStringFromReqBody(req))

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			fmt.Println(computeErr)
			return
		}

		if len(data.loadedKey) > 0 {
			fmt.Println("Loaded key: " + data.loadedKey)
			CallScript(&data.interpretor, &inputPath, &outputPath, &data.loadedKey, "decrypt")
		} else {
			fmt.Println("Cannot decrypt because no key has been found")
		}
	}

	var handleAddKeyAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		var inputKeyPath string = getStringFromReqBody(req)
		inputKeyPath = strings.Replace(inputKeyPath, "\"", "", -1)
		outputKeyPath := GetKeysDirPath() + "/" + keymgn.GenerateKeyName()

		fmt.Println("Add key action was triggered: " + inputKeyPath)

		data.loadedKey = keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
	}

	// Endpoints handlers
	http.HandleFunc("/encrypt", handleEncryptAction)
	http.HandleFunc("/decrypt", handleDecryptAction)
	http.HandleFunc("/addkey", handleAddKeyAction)
}

func SetupAppDirs() error {
	//TODO the creation part of the directories is going to stay in the installer
	if _, err := os.Stat(GetAppDirPath()); errors.Is(err, os.ErrNotExist) {
		if createDirErr := os.Mkdir(GetAppDirPath(), os.ModePerm); createDirErr != nil {
			log.Println("Cannot create app directory: ", createDirErr)
			return err
		}
	}

	if _, err := os.Stat(GetBinDirPath()); errors.Is(err, os.ErrNotExist) {
		if createDirErr := os.Mkdir(GetBinDirPath(), os.ModePerm); createDirErr != nil {
			log.Println("Cannot create client directory: ", createDirErr)
			return err
		}
	}

	if _, err := os.Stat(GetKeysDirPath()); errors.Is(err, os.ErrNotExist) {
		if createDirErr := os.Mkdir(GetKeysDirPath(), os.ModePerm); createDirErr != nil {
			log.Println("Cannot create keys directory: ", createDirErr)
			return err
		}
	}

	return nil
}

func Run() {
	PORT := "1234"
	fmt.Printf("Server has been started on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func CallScript(pythonExecPath *string, inputPath *string, outputPath *string, loadedKey *string, action string) {
	scriptPath := filepath.Join(GetBinDirPath(), "filecrypt.py")
	c := exec.Command(*pythonExecPath, scriptPath, action, *loadedKey, *inputPath, *outputPath)

	if out, err := c.Output(); err != nil {
		fmt.Println("Error when encrypting: ", err)
		fmt.Println("Command output: ", out)
	}
}

func ComputeOutputPath(inputPath *string) (string, error) {
	inputInfo, err := os.Stat(*inputPath)
	if err != nil {
		return "", errors.New("Cannot get the stats of path: " + *inputPath)
	}

	var outputPath string = filepath.Dir(*inputPath)

	if inputInfo.IsDir() {
		filename := filepath.Base(*inputPath)
		outputPath = outputPath + "/" + filename + "_result"
		if createDirErr := os.Mkdir(outputPath, os.ModePerm); createDirErr != nil {
			return "", errors.New("Cannot create directory in: " + outputPath)
		}
	}

	return outputPath, nil
}

func getStringFromReqBody(req *http.Request) string {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
