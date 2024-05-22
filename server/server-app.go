package main

import (
	"server/env"
	"server/utils"
)

func main() {
	log := utils.GetLogger()
	log.Info("Starting server...")

	var env *env.Environment = new(env.Environment)

	env.Setup()
	env.Run()

	log.Info("Server is shutting down")
}
