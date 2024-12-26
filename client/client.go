package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var username = ""
var scanner = bufio.NewReader(os.Stdin)

func main() {
	fmt.Print("Username:\t")
	username, _ = scanner.ReadString('\n')
	username = strings.TrimSuffix(username, "\r\n")

	sendData()
	fmt.Println("Goodbye!")
}

func sendData() {

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
		message = strings.TrimSuffix(message, "\r\n")
		if message == "STOP" {
			break
		}

		newMessage := makeMessage(message)
		conn.Write(newMessage)
	}
}

func makeMessage(content string) []byte {
	message := make([]byte, 0)
	message = append(message, pad([]byte("main"), 32)...)
	message = append(message, pad([]byte("message"), 32)...)
	message = append(message, pad([]byte(username), 64)...)
	message = append(message, pad([]byte(content), 128)...)
	return message
}

func pad(b []byte, size int) []byte {
	l := len(b)
	if l == size {
		return b
	}
	if l > size {
		return b[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], b)
	return tmp
}
