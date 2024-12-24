package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	sendData()
	fmt.Println("Goodbye!")
}

func sendData() {
	scanner := bufio.NewReader(os.Stdin)

	conn, connErr := net.Dial("tcp", "localhost:8001")
	if connErr != nil {
		fmt.Println(connErr)
		return
	}

	defer conn.Close()
	go func() {
		for {
			reply := make([]byte, 1024)
			_, err := conn.Read(reply)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(reply))
		}
	}()
	for {
		message, _ := scanner.ReadString('\n')

		if message == "STOP" {
			break
		}

		conn.Write([]byte(message + "\n"))
	}
}
