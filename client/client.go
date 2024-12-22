package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	sendData()
}

func sendData() {
	scanner := bufio.NewReader(os.Stdin)

	conn, connErr := net.Dial("tcp", "localhost:8001")
	if connErr != nil {
		fmt.Println(connErr)
		return
	}

	defer conn.Close()

	for {
		reply := make([]byte, 1024)
		_, err := conn.Read(reply)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(reply))

		message, _ := scanner.ReadString('\n')

		conn.Write([]byte(message + "\n"))
	}
	fmt.Println("Goodbye")
}
