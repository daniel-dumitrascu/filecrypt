package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Exiting - path argument wasn't provided")
		return
	}

	c, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	pathString := arguments[1]

	fmt.Fprintf(c, pathString+"\n")
}
