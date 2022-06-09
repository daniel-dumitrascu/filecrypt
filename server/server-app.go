package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {

	/*if ok, msg := env.Setup(); !ok {
		fmt.Println("There was an issue in setting up the env: " + msg)
		return
	}

	server.Start()*/

	http.HandleFunc("/encrypt", handleEncryptAction)
	http.HandleFunc("/decrypt", handleDecryptAction)

	PORT := "1234"
	http.ListenAndServe(":"+PORT, nil)
}
