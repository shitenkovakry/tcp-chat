package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
)

const (
	address            = "localhost:8080"
	stopWordFromClient = "goodbye"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()

	addressOfClient := connection.RemoteAddr().String()
	log.Print("connected with client:", addressOfClient)

	bufferWindow := make([]byte, 1024)

	readLen, err := connection.Read(bufferWindow)
	if err != nil {
		log.Print(errors.Wrapf(err, "can not read message from client"))

		return
	}

	nicknameFromClient := string(bufferWindow[:readLen])
	log.Println("received nickname from client. nickname:", nicknameFromClient)

	response := "nickname received and read"

	_, err = connection.Write([]byte(response))
	if err != nil {
		log.Println(errors.Wrapf(err, "can not send response to client"))

		return
	}

	log.Println("response successfully sent to client")
}

func startServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println(errors.Wrapf(err, "can not listen connection"))
	}
	defer listener.Close()

	fmt.Println("server listening at:", address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(errors.Wrapf(err, "can not accept connection"))

			continue
		}

		go handleConnection(connection)
	}
}

func main() {
	startServer(address)
}
