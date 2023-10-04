package main

import (
	"fmt"
	"server/env"
)

func main() {
	fmt.Printf("Starting server...\n")
	/*if ok, msg := env.Setup(); !ok {
		fmt.Println("There was an issue in setting up the env: " + msg)
		return
	}

	server.Start()*/

	var data env.EnvData

	env.Setup(&data)
	env.Run()
	fmt.Printf("Starting is shutting down\n")
}
