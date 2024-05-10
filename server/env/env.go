package env

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/crypt"
	"server/keymgn"
	"server/request"
	"server/utils"
	"server/workerpool"
)

type Environment struct {
	loadedKey []byte
	pool      workerpool.Pool
}

func (env *Environment) Setup() {
	log := utils.GetLogger()
	log.Info("Setup the environment")
	osmanager := GetOsManager()
	osmanager.SpecificSetup()
	var installKeyPath = osmanager.GetKeysDirPath()
	env.loadedKey = keymgn.LoadKey(&installKeyPath)

	var handleEncryptAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputPath := req.TargetPath
		log.Info("Encrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			log.Error(computeErr)
			return
		}

		if len(env.loadedKey) > 0 {
			crypt.EncryptDir(inputPath, outputPath, env.loadedKey)
			log.Info("Successfully encrypting the file: " + inputPath)
		} else {
			log.Error("Cannot encrypt because no key has been found")
		}
	}

	var handleDecryptAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputPath := req.TargetPath
		log.Info("Decrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			log.Error(computeErr)
			return
		}

		if len(env.loadedKey) > 0 {
			crypt.DecryptDir(inputPath, outputPath, env.loadedKey)
			log.Info("Successfully decrypting the file: " + inputPath)
		} else {
			log.Error("Cannot decrypt because no key has been found")
		}
	}

	var handleAddKeyAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputKeyPath := req.TargetPath
		outputKeyPath := osmanager.GetKeysDirPath() + "/" + keymgn.GenerateKeyName()

		log.Info("Add key action was triggered: " + inputKeyPath)

		keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
		env.loadedKey = keymgn.LoadKey(&outputKeyPath)
		osmanager.ChangeFilePermission(&outputKeyPath)
	}

	var handleGenKeyAction func(req *request.RequestData) = func(req *request.RequestData) {
		keyname := keymgn.GenerateKeyName()
		outputKeyPath := req.TargetPath + "/" + keyname

		log.Info("Gen key action was triggered: " + outputKeyPath)
		key := crypt.GenKey()
		encodedKey := base64.StdEncoding.EncodeToString(key)

		if err := os.WriteFile(keyname, []byte(encodedKey), 0644); err != nil {
			log.Error("Cannot save key to file")
			return
		}
	}

	var handlers [4]func(req *request.RequestData)
	handlers[0] = handleEncryptAction
	handlers[1] = handleDecryptAction
	handlers[2] = handleAddKeyAction
	handlers[3] = handleGenKeyAction

	env.pool.Init(config.Max_goroutines_nr, &handlers)

	// Endpoints handlers
	http.HandleFunc("/process", env.processHandler)
}

func (env *Environment) Run() {
	PORT := "1234"
	log := utils.GetLogger()
	log.Info("Server has been started on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func (env *Environment) processHandler(w http.ResponseWriter, req *http.Request) {
	var reqData request.RequestData
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log := utils.GetLogger()

	if reqData.ActionType < 0 || reqData.ActionType > config.Max_handlers_nr {
		log.Info("The action type %d is invalid. This should be between 0 and %d",
			reqData.ActionType, config.Max_handlers_nr-1)
		return
	}

	env.pool.AddTask(&reqData)
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
