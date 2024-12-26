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
	"sync"
)

type Account struct {
	username  string
	encrypted string
}

type Message struct {
	content  string
	from     string
	roomName string
	intent   string
}

type Room struct {
	mu       sync.Mutex
	channels []net.Conn
}

type Server struct {
	mu    sync.Mutex
	rooms map[string]Room
}

var globalServer Server

func (r *Room) broadcast(message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.channels {
		c.Write([]byte(message + "\n"))
	}
}

func (r *Room) add(c net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels = append(r.channels, c)
}

func (r *Room) remove(c net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, v := range r.channels {
		if v == c {
			r.channels[i] = r.channels[len(r.channels)-1]
			r.channels = r.channels[:len(r.channels)-1]
			return
		}
	}
}

func (r *Room) printMembers() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.channels {
		fmt.Println(c.RemoteAddr().String())
	}
}

func (s *Server) addRoom(roomName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms[roomName] = Room{}
}

func (s *Server) removeRoom(roomName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rooms, roomName)
}

func (s *Server) listRooms() (r []string) {
	for k, _ := range s.rooms {
		r = append(r, k)
	}
	return
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
	mainRoom := Room{
		channels: []net.Conn{},
	}
	l, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer l.Close()

	fmt.Println("Serving on", PORT)

	for {
		c, err := l.Accept()
		fmt.Println("Accepted connection")
		mainRoom.add(c)
		fmt.Println("Added connection to mainroom")

		fmt.Println("Current room status:")
		mainRoom.printMembers()

		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleMessageConnection(c, &mainRoom)
	}
}

func handleMessageConnection(c net.Conn, r *Room) {
	defer c.Close()
	defer r.remove(c)
	defer r.broadcast(c.RemoteAddr().String() + " has left the room")

	fmt.Println("New connection from", c.RemoteAddr())

	r.broadcast("Welcome! " + c.RemoteAddr().String())

	for {
		netDataBytes, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			log.Println("Connection closed or error:", err)
			return
		}
		netData := string(netDataBytes)

		fmt.Println(netData)

		tempMessage := receive(netDataBytes)
		if tempMessage.intent == "STOP" {
			fmt.Println("Stopping connection with", c.RemoteAddr())
			break
		}

		fmt.Println(tempMessage.content)

		r.broadcast(tempMessage.content)
	}
}

func receive(bytesData []byte) (m Message) {
	fmt.Println(bytesData)
	fmt.Println(len(bytesData))

	m.roomName = string(trimPadding(bytesData[:32], 32))
	m.intent = string(trimPadding(bytesData[32:64], 32))
	m.from = string(trimPadding(bytesData[64:128], 64))
	m.content = string(trimPadding(bytesData[128:256], 128))

	fmt.Printf("%+v\n", m)

	return m
}

func trimPadding(bytesData []byte, size int) []byte {
	return bytesData[len(bytesData)-size:]
}
