package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	address  = "localhost:8080"
	stopWord = "goodbye"
)

func main() {
	// Создание канала для обработки сигнала
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	// Подключение к серверу по адресу "localhost:8080"
	connection, err := net.Dial("tcp", address)
	if err != nil {
		panic(errors.Wrapf(err, "can not connection to server"))
	}

	defer connection.Close()

	nickName := "kry"

	messages := []string{
		"durik",
		"ondrys",
		// "miu",
		// "",
	}

	// Чтение данных от сервера
	for index := 0; index < len(messages); index++ {
		time.Sleep(time.Second * 10)

		message := []byte(messages[index])

		wroteBytes, err := connection.Write(message)
		if err != nil {
			log.Print(errors.Wrapf(err, "can not write data"))
			return
		}

		log.Print("wrote ", wroteBytes, " bytes")
	}

	// Ожидание сигнала завершения программы
	<-signalChan

	_, err = connection.Write([]byte(stopWord))
	if err != nil {
		log.Print(errors.Wrapf(err, "can not send information to server"))
		return
	}

	log.Println("the word", stopWord, "was sent to server. connection will be closed")

	log.Println("the programm ends. close the connection")
}
