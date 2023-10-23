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

	var env *env.Environment = new(env.Environment)

	env.Setup()
	env.Run()

	fmt.Printf("Server is shutting down\n")
}
