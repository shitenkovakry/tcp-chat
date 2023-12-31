package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

const (
	address            = "localhost:9998"
	addressOfCompanion = "localhost:9999"
	stopWord           = "goodbye"
)

func sendToServerNickname(connection net.Conn, nickname string) error {
	_, err := connection.Write([]byte(nickname))
	if err != nil {
		return errors.Wrapf(err, "can not send name to server")
	}

	fmt.Println("nickname", nickname, "successfully sent to server")

	return nil
}

func sendRequestToServer(connection net.Conn, addressOfCompanion string, nicknameOfCompanion string) error {
	request := fmt.Sprintln("connect with", nicknameOfCompanion)

	if _, err := connection.Write([]byte(request)); err != nil {
		return errors.Wrapf(err, "can not connect")
	}

	log.Println("communication request sent", nicknameOfCompanion)

	// получение ответа от сервера
	response := make([]byte, 1024)

	readLen, err := connection.Read(response)
	if err != nil {
		return errors.Wrapf(err, "can not read response from server")
	}

	log.Println("recieved a response from server:", string(response[:readLen]))

	return nil
}

func connectToServer(address string) (net.Conn, error) {
	// Подключение к серверу по адресу "localhost:8080"
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, errors.Wrapf(err, "can not connection to server")
	}

	return connection, nil
}

func finishServer(connection net.Conn) {
	defer connection.Close()

	if _, err := connection.Write([]byte(stopWord)); err != nil {
		log.Println(errors.Wrapf(err, "can not send information to server"))

		return

	}

	log.Println("the word", stopWord, "was sent to server. connection will be closed")

	log.Println("the programm ends. close the connection")
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	nickname := "kry"
	nicknameOfCompanion := "ondrys"

	connection, err := connectToServer(address)
	if err != nil {
		panic(err)
	}

	defer finishServer(connection)

	if err := sendToServerNickname(connection, nickname); err != nil {
		panic(errors.Wrapf(err, "can not send name to server"))

	}

	if err := sendRequestToServer(connection, addressOfCompanion, nicknameOfCompanion); err != nil {
		panic(errors.Wrapf(err, "can not send request to companion"))
	}

	// Ожидание сигнала завершения программы
	<-signalChan

}
