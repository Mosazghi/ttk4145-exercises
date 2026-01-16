package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"
)

/**
* Network functionality
* - connect to UDP or TCP
* - Write
* - Recieve
*
* Basic functionality:
* - write
* */

type DataDirection int

const (
	SEND = iota
	RECEIVE
)

type NetworkInfo struct {
	RemoteAddress string
	SendPort      string
	ReceivePort   string
	ProtocolType  string
}

func Connect(network *NetworkInfo, direction DataDirection) net.Conn {
	address := network.RemoteAddress

	if direction == SEND {
		address = address + ":" + network.SendPort
	}

	if direction == RECEIVE {
		address = address + ":" + network.ReceivePort
	}

	conn, err := net.Dial(network.ProtocolType, address)
	if err != nil {
		fmt.Println(err)
	}

	return conn
}

func GetMessage() (string, error) {
	var message string
	fmt.Print("Input: ")
	_, err := fmt.Scanln(&message)
	if err != nil {
		fmt.Println("Failed to get message, got err ", err)
		return "", err
	}

	return message, err
}

func Write(connection net.Conn, message string) {
	if _, err := connection.Write([]byte(message)); err != nil {
		fmt.Println("failed to write to server, got err ", err)
		return
	}

	fmt.Println("Wrote successfully to server, message: ", message)
}

func Read(connection net.Conn) {
	message := make([]byte, 1024)
	if _, err := connection.Read(message); err != nil {
		fmt.Println("failed to read from server, got err ", err)
		return
	}

	fmt.Println("Got from server: ", string(message))
}

func CreateUI(protocol string) {
	options := []string{"Send message", "Read message", "Quit"}
	fmt.Println("-------------------------------------")
	fmt.Println("             ", protocol, "DEMO ")
	fmt.Println("-------------------------------------")

	for index, option := range options {
		fmt.Println(index+1, ": ", option)
	}
}

func CheckInput(input string) bool {
	legalInputs := []string{"1", "2", "3"}

	if !slices.Contains(legalInputs, input) {
		fmt.Println("Not legal input try again...")
		return false
	}
	return true
}

func client(networkInfo *NetworkInfo) {
	recieverConnection := Connect(networkInfo, RECEIVE)
	senderConnection := Connect(networkInfo, SEND)
	running := true

	for running {
		CreateUI(networkInfo.ProtocolType)

		input, err := GetMessage()
		if err != nil {
			continue
		}

		if res := CheckInput(input); !res {
			continue
		}

		switch input {
		case "1":
			message, err := GetMessage()
			if err != nil {
				continue
			}
			Write(senderConnection, message)

		case "2":
			Read(recieverConnection)

		case "3":
			running = false
		}
	}
}

func main() {
	args := os.Args
	allowedProtocols := []string{"udp", "tcp"}

	if len(args) < 3 {
		log.Fatal("Too few arguments")
	}

	if !slices.Contains(allowedProtocols, args[1]) {
		log.Fatal("Not allowed network type, got ", args[1])
	}

	if net.ParseIP(args[2]) == nil {
		log.Fatal("Wrong formated IP address, got ", args[2])
	}

	netInfo := NetworkInfo{
		args[2],
		"20000",
		"30000",
		args[1],
	}

	client(&netInfo)
}
