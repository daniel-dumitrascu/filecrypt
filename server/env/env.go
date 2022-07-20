package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server_app/keymgn"
	"strings"
)

type EnvData struct {
	loadedKey string
}

func Setup(data *EnvData) {
	osmanager := GetOsManager()

	if _, err := os.Stat(osmanager.GetContextAppPath()); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Context app is not installed.")
	}

	osmanager.SpecificSetup()
	var installKeyPath = osmanager.GetInstallKeyPath()
	data.loadedKey = keymgn.LoadKey(&installKeyPath)

	var handleEncryptAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Encrypt action was triggered: " + getStringFromReqBody(req))
		fmt.Println("Loaded key: " + data.loadedKey)
	}

	var handleDecryptAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Decrypt action was triggered: " + getStringFromReqBody(req))
		fmt.Println("Loaded key: " + data.loadedKey)
	}

	var handleAddKeyAction func(w http.ResponseWriter, req *http.Request) = func(w http.ResponseWriter, req *http.Request) {
		var inputKeyPath string = getStringFromReqBody(req)
		inputKeyPath = strings.Replace(inputKeyPath, "\"", "", -1)
		outputKeyPath := GetInstallKeyPath() + keymgn.GenerateKeyName()

		fmt.Println("Add key action was triggered: " + inputKeyPath)

		keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
	}

	// Endpoints handlers
	http.HandleFunc("/encrypt", handleEncryptAction)
	http.HandleFunc("/decrypt", handleDecryptAction)
	http.HandleFunc("/addkey", handleAddKeyAction)
}

func Run() {
	PORT := "1234"
	http.ListenAndServe(":"+PORT, nil)
}

func GetInstallKeyPath() string {
	osmanager := GetOsManager()
	return osmanager.GetInstallKeyPath()
}

func getStringFromReqBody(req *http.Request) string {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
