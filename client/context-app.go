package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func computeAction(targetPath *string) string {
	start := strings.LastIndex(*targetPath, ".")

	if start != -1 {
		extension := (*targetPath)[start+1:]
		if extension == "crypt" {
			return "decrypt"
		}
	}

	return "encrypt"
}

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("Exiting - action and path argument weren't provided")
		return
	}

	action := arguments[1]
	targetPath := string(arguments[2])

	var endpoint string
	if action == "crypt" {
		endpoint = computeAction(&targetPath)
	} else if action == "addkey" {
		endpoint = "addkey"
	} else {
		log.Fatalf("Context-app has been called without a valid action")
	}

	postBody, _ := json.Marshal(targetPath)
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
