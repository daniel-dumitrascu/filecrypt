package main

import (
	"server_app/env"
)

func main() {

	/*if ok, msg := env.Setup(); !ok {
		fmt.Println("There was an issue in setting up the env: " + msg)
		return
	}

	server.Start()*/

	var data env.EnvData

	env.Setup(&data)
	env.Run()
}
