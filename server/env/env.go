package env

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/config"
	"server/keymgn"
	"server/request"
	"server/utils"
	"server/workerpool"
)

type Environment struct {
	key  string
	pool workerpool.Pool
}

func (env *Environment) Setup() {
	log := utils.GetLogger()
	log.Info("Setup the environment")
	osmanager := GetOsManager()
	osmanager.SpecificSetup()
	var installKeyPath = osmanager.GetKeysDirPath()
	env.key = keymgn.GetLatestKey(&installKeyPath)
	toolPath := osmanager.GetBinDirPath() + "/" + config.CRYPT_TOOL_NAME

	if config.CURRENT_PLATFORM == config.PLATFORM_WIN {
		toolPath += ".exe"
	}

	var handleEncryptAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputPath := req.TargetPath
		log.Info("Encrypt action was triggered for: " + inputPath)

		outputPath, computeErr := ComputeOutputPath(&inputPath)
		if computeErr != nil {
			log.Error(computeErr)
			return
		}

		if len(env.key) > 0 {
			callToolEncrypt(&toolPath, &inputPath, &outputPath, &env.key)
			log.Info("Encryption task was completed successfully!")
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

		if len(env.key) > 0 {
			callToolDecrypt(&toolPath, &inputPath, &outputPath, &env.key)
			log.Info("Decryption task was completed successfully!")
		} else {
			log.Error("Cannot decrypt because no key has been found")
		}
	}

	var handleAddKeyAction func(req *request.RequestData) = func(req *request.RequestData) {
		inputKeyPath := req.TargetPath
		if isdir, _ := utils.IsDir(inputKeyPath); isdir {
			log.Error("Target path is a directory.")
			return
		}

		outputKeyPath := osmanager.GetKeysDirPath() + "/" + keymgn.GenerateKeyName()

		log.Info("Add key action was triggered: " + inputKeyPath)

		keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
		env.key = keymgn.GetLatestKey(&outputKeyPath)
		osmanager.ChangeFilePermission(&outputKeyPath)
	}

	var handleGenKeyAction func(req *request.RequestData) = func(req *request.RequestData) {
		outputKeyPath := req.TargetPath
		log.Info("Initial path: ", outputKeyPath)
		if isdir, _ := utils.IsDir(outputKeyPath); !isdir {
			outputKeyPath = filepath.Dir(outputKeyPath)
		}

		keyname := keymgn.GenerateKeyName()
		outputKeyPath = outputKeyPath + "/" + keyname

		log.Info("Calculated path: ", outputKeyPath)

		log.Info("Gen key action was triggered: " + outputKeyPath)
		key := callToolGenKey(&toolPath)
		if key == nil {
			return
		}

		encodedKey := base64.StdEncoding.EncodeToString(key)
		if err := os.WriteFile(outputKeyPath, []byte(encodedKey), 0644); err != nil {
			log.Error("Cannot save key to file")
			return
		}
	}

	var handlers [4]func(req *request.RequestData)
	handlers[0] = handleEncryptAction
	handlers[1] = handleDecryptAction
	handlers[2] = handleAddKeyAction
	handlers[3] = handleGenKeyAction

	env.pool.Init(config.MAX_GOROUTINES_NR, &handlers)

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

	if reqData.ActionType < 0 || reqData.ActionType > config.MAX_HANDLERS_NR {
		log.Info("The action type %d is invalid. This should be between 0 and %d",
			reqData.ActionType, config.MAX_HANDLERS_NR-1)
		return
	}

	env.pool.AddTask(&reqData)
}

func callToolGenKey(toolPath *string) []byte {
	c := exec.Command(*toolPath, "genkey")

	log := utils.GetLogger()
	out, err := c.Output()

	if err != nil {
		log.Error("Error when generating a new key: ", err)
		log.Error("Command output: ", out)
		return nil
	}

	return out
}

func callToolDecrypt(toolPath *string, inputPath *string, outputPath *string, keyPath *string) {
	callCryptTool(toolPath, inputPath, outputPath, keyPath, "decrypt")
}

func callToolEncrypt(toolPath *string, inputPath *string, outputPath *string, keyPath *string) {
	callCryptTool(toolPath, inputPath, outputPath, keyPath, "encrypt")
}

func callCryptTool(toolPath *string, inputPath *string, outputPath *string, keyPath *string, action string) {
	c := exec.Command(*toolPath, action, *keyPath, *inputPath, *outputPath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	log := utils.GetLogger()

	if err := c.Run(); err != nil {
		log.Error("Error when encrypting: ", err)
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
