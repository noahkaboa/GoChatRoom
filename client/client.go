package main

import (
	"fmt"
	"net"
)

func main() {
	sendData()
}

func sendData() {
	conn, connErr := net.Dial("tcp", "localhost:8001")
	if connErr != nil {
		fmt.Println(connErr)
		return
	}

	defer conn.Close()

	for {
		var message string
		fmt.Print(">\t")
		fmt.Scan(&message)
		fmt.Println("|" + message + "|")

		conn.Write([]byte(message + "\n"))

		reply := make([]byte, 1024)
		_, err := conn.Read(reply)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(reply))
	}
	fmt.Println("Goodbye")

}
