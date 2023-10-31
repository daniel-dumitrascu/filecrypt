package main

import (
	"server/env"
	"server/utils"
)

func main() {
	log := utils.GetLogger()
	log.Info("Starting server...")
	/*if ok, msg := env.Setup(); !ok {
		fmt.Println("There was an issue in setting up the env: " + msg)
		return
	}

	server.Start()*/

	var env *env.Environment = new(env.Environment)

	env.Setup()
	env.Run()

	log.Info("Server is shutting down")
}
