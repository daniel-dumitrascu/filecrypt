package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	PORT := ":1234"
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Message received -> ", string(netData))
	}
}
