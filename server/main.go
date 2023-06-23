package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/pkg/errors"
)

const (
	address            = "localhost:9998"
	addressOfCompanion = "localhost:9999"
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

	nicknameOfCompanion, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		log.Println(errors.Wrapf(err, "can not read nickname of another client:", nicknameOfCompanion))

		return
	}

	nicknameOfCompanion = strings.TrimSpace(nicknameFromClient)
	log.Println("received nickname from client:", nicknameOfCompanion)

	// отправка запроса другому клиенту
	request := fmt.Sprintln("request from:", nicknameOfCompanion)

	response := "nickname received and read"
	responseFromCompanion, err := sendRequestToServerOfCompanion(request, addressOfCompanion)
	if err != nil {
		log.Println(errors.Wrapf(err, "can not send request to another client"))
	}

	log.Println("you have answer from other client:", response)

	_, err = connection.Write([]byte(response))
	if err != nil {
		log.Println(errors.Wrapf(err, "can not send response to client"))

		return
	}

	// отправление ответа клиенту, который запросил общение
	_, err = connection.Write([]byte(responseFromCompanion + "\n"))
	if err != nil {
		log.Println(errors.Wrapf(err, "error while sending response fom anothe client"))

		return
	}

	log.Println("response successfully sent to client")
	log.Println("reply sent from another client", response)
}

func sendRequestToServerOfCompanion(request string, address string) (string, error) {
	addressOfCompanion := address

	connectionWithCompanion, err := net.Dial("tcp", addressOfCompanion)
	if err != nil {
		return "", errors.Wrapf(err, "connection can not br established ")
	}
	defer connectionWithCompanion.Close()

	_, err = connectionWithCompanion.Write([]byte(request))
	if err != nil {
		return "", errors.Wrapf(err, "can not connect to another client")
	}

	log.Println("request sent to another client", request)

	// чтение ответа другого клиента
	response, err := bufio.NewReader(connectionWithCompanion).ReadString('\n')
	if err != nil {
		return "", errors.Wrapf(err, "error while reading another respons of client")
	}

	response = strings.TrimSpace(response)

	return response, nil
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
