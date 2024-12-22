package main

import (
	"bufio"
	"bytes"
	b64 "encoding/base64"
	csv "encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Account struct {
	username  string
	encrypted string
}

type Message struct {
	content string
	from    string
	to      string
}

const databasePath = "db.csv"
const PORT = ":8001"

func main() {
	// accounts, readErr := getDB()
	// fmt.Println(readErr)
	// fmt.Println(accounts)
	// writeErr := newAccount()
	// fmt.Println(writeErr)
	// accounts, readErr = getDB()
	// fmt.Println(readErr)
	// fmt.Println(accounts)

	serve()
	fmt.Println("Done serving")

}

func newAccount() error {
	var username string
	var password string
	var encodedPassword string

	fmt.Println("What is the username?")
	fmt.Print(">\t")
	fmt.Scan(&username)

	fmt.Println("What is the password?")
	fmt.Print(">\t")
	fmt.Scan(&password)

	encodedPassword = b64.StdEncoding.EncodeToString([]byte(password))

	dbWriteErr := writeDB([]string{username, encodedPassword})

	return dbWriteErr
}

func getDB() ([]Account, error) {
	file, err := os.Open(databasePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(bytes.NewReader(data))

	var accounts []Account

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return nil, err //break?
		}
		a := Account{
			username:  record[0],
			encrypted: record[1],
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

func writeDB(record []string) error {
	f, fileErr := os.OpenFile(databasePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if fileErr != nil {
		return fileErr
	}
	writer := csv.NewWriter(f)
	writingErr := writer.Write(record)

	writer.Flush()

	return writingErr
}

func serve() {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer l.Close()

	fmt.Println("Serving on", PORT)

	for {
		c, err := l.Accept()
		fmt.Println("Accepted connection")
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleMessageConnection(c)
	}
}

func handleMessageConnection(c net.Conn) {
	defer c.Close()
	fmt.Println("New connection from", c.RemoteAddr())

	_, err := c.Write([]byte("Welcome!"))
	if err != nil {
		fmt.Println(err)
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println("Connection closed or error:", err)
			return
		}

		temp := strings.TrimSpace(netData)
		if temp == "STOP" {
			fmt.Println("Stopping connection with", c.RemoteAddr())
			break
		}

		fmt.Println(temp)

		_, err = c.Write([]byte(temp))
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}
	}
}
