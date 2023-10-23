package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/config"
	"server/keymgn"
	"server/request"
	"server/workerpool"
)

type Environment struct {
	loadedKey   string
	interpretor string
	pool        workerpool.Pool
}

func (env *Environment) Setup() {
	fmt.Printf("Setup the environment\n")
	if err := setupAppDirs(); errors.Is(err, os.ErrNotExist) {
		os.Exit(1)
	}

	osmanager := GetOsManager()
	osmanager.SpecificSetup()
	var installKeyPath = GetKeysDirPath()
	env.loadedKey = keymgn.LoadKey(&installKeyPath)
	env.interpretor = osmanager.GetInterpretor()

	var handleEncryptAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputPath := req.TargetPath
		fmt.Println("Encrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			fmt.Println(computeErr)
			return
		}

		if len(env.loadedKey) > 0 {
			fmt.Println("Loaded key: " + env.loadedKey)
			CallScript(&env.interpretor, &inputPath, &outputPath, &env.loadedKey, "encrypt")
		} else {
			fmt.Println("Cannot encrypt because no key has been found")
		}
	}

	var handleDecryptAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputPath := req.TargetPath
		fmt.Println("Decrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			fmt.Println(computeErr)
			return
		}

		if len(env.loadedKey) > 0 {
			fmt.Println("Loaded key: " + env.loadedKey)
			CallScript(&env.interpretor, &inputPath, &outputPath, &env.loadedKey, "decrypt")
		} else {
			fmt.Println("Cannot decrypt because no key has been found")
		}
	}

	var handleAddKeyAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputKeyPath := req.TargetPath
		outputKeyPath := GetKeysDirPath() + "/" + keymgn.GenerateKeyName()

		fmt.Println("Add key action was triggered: " + inputKeyPath)

		env.loadedKey = keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
	}

	var handlers [3]func(req *request.RequestData)
	handlers[0] = handleEncryptAction
	handlers[1] = handleDecryptAction
	handlers[2] = handleAddKeyAction

	env.pool.Init(config.Max_goroutines_nr, &handlers)

	// Endpoints handlers
	http.HandleFunc("/process", env.processHandler)
}

func (env *Environment) Run() {
	PORT := "1234"
	fmt.Printf("Server has been started on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func (env *Environment) processHandler(w http.ResponseWriter, req *http.Request) {
	var reqData request.RequestData
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqData.ActionType < 0 || reqData.ActionType > config.Max_handlers_nr {
		fmt.Printf("The action type %d is invalid. This should be between 0 and %d",
			reqData.ActionType, config.Max_handlers_nr-1)
		return
	}

	env.pool.AddTask(&reqData)
}

func setupAppDirs() error {
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

// TODO to remove
func getStringFromReqBody(req *http.Request) string {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
