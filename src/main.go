package main

import (
	"log"
	"net"
)

const (
	ImageFolder  = "/app/images"
	InputFolder  = "/app/input"
	SharedFolder = "/app/shared"
	OutputFolder = "/app/output"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// read data from the connection
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	// write data back to the connection
	conn.Write([]byte("Hello, client!"))

	// close the connection
	conn.Close()
}
