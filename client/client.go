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
	f, _ := os.Create("C:\\Users\\DanielDumitrascu\\Desktop\\log\\log.txt")
	f.Write([]byte("in main()" + "\n"))
	arguments := os.Args
	//arguments := []string{"", "encrypt", "C:\\Users\\DanielDumitrascu\\Desktop\\1"}

	if len(arguments) < 3 {
		log.Fatalf("Exiting - action and path argument weren't provided")
	}

	action := arguments[1]
	targetPath := string(arguments[2] + "\\")

	f.Write([]byte("Action: " + action + "\n"))
	f.Write([]byte("Path: " + targetPath + "\n"))
	f.Write([]byte("---------------------------------" + "\n"))
	f.Close()

	var reqData RequestData
	if action == "encrypt" {
		reqData.ActionType = 0
	} else if action == "decrypt" {
		reqData.ActionType = 1
	} else if action == "addkey" {
		reqData.ActionType = 2
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
