package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("Exiting - action and path argument weren't provided")
		return
	}

	action := arguments[1]
	var endpoint string
	if action == "crypt" {
		endpoint = "encrypt"
	} else if action == "addkey" {
		endpoint = "addkey"
	} else {
		log.Fatalf("Context-app has been called without a valid action")
	}

	postBody, _ := json.Marshal(string(arguments[2]))
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://127.0.0.1:1234/"+endpoint, "application/json", responseBody)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
