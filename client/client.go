package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		log.Fatalf("Exiting - action and path argument weren't provided")
	}

	action := arguments[1]
	targetPath := string(arguments[2])

	var reqData RequestData
	if action == "encrypt" {
		reqData.ActionType = 0
	} else if action == "decrypt" {
		reqData.ActionType = 1
	} else if action == "addkey" {
		reqData.ActionType = 2
	} else if action == "genkey" {
		reqData.ActionType = 3
	} else {
		log.Fatalf("Context-app has been called without a valid action")
	}

	reqData.TargetPath = targetPath

	postBody, _ := json.Marshal(reqData)
	requestBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(HomeAddress+":"+HomePort+"/process", MediaType, requestBody)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Println(sb)
}
