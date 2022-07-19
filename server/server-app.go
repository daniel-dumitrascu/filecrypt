package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server_app/env"
	"server_app/keymgn"
	"strings"
)

func getStringFromReqBody(req *http.Request) string {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func handleEncryptAction(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Encrypt action was triggered: " + getStringFromReqBody(req))
}

func handleDecryptAction(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Decrypt action was triggered: " + getStringFromReqBody(req))
}

func handleAddKeyAction(w http.ResponseWriter, req *http.Request) {
	var inputKeyPath string = getStringFromReqBody(req)
	inputKeyPath = strings.Replace(inputKeyPath, "\"", "", -1)
	outputKeyPath := env.GetInstallKeyPath() + keymgn.GenerateKeyName()

	fmt.Println("Add key action was triggered: " + inputKeyPath)

	keymgn.InstallKey(&inputKeyPath, &outputKeyPath)
}

func main() {

	/*if ok, msg := env.Setup(); !ok {
		fmt.Println("There was an issue in setting up the env: " + msg)
		return
	}

	server.Start()*/

	env.Setup()

	http.HandleFunc("/encrypt", handleEncryptAction)
	http.HandleFunc("/decrypt", handleDecryptAction)
	http.HandleFunc("/addkey", handleAddKeyAction)

	PORT := "1234"
	http.ListenAndServe(":"+PORT, nil)
}
