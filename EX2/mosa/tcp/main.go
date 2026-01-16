package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const SERVER_IP = "10.100.23.11"
const SERVER_PORT = "34933"

func tcpListener() (net.Conn, error) {
	fmt.Println("Listening...")
	addr, err := net.ResolveTCPAddr("tcp", SERVER_IP+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("[ECHO] Error resolving addr: ", err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("[ECHO] Error listening: ", err)
		return nil, err
	}
	return conn, nil
}

func tcpReceiver(conn net.Conn, recvChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Err reading...", err)
			return
		}
		data := string(buffer[:n])
		recvChan <- data

	}
}

func tcpSender(conn net.Conn, inputChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	msgReader := bufio.NewReader(os.Stdin)
	fmt.Println("[SENDER] Connected to ", conn.LocalAddr())

	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		// fmt.Print("> ")
		n, err := msgReader.Read(buffer)
		// message = strings.TrimSuffix(message, "\n")
		if err != nil {
			fmt.Println("Error reading message: ", err)
			continue
		}
		if n == 1 {
			continue
		}

		_, err = conn.Write(append(buffer[:n], 0))
		if err != nil {
			fmt.Println("ERROR: Write failed")
		}

		inputChan <- string(buffer[:n])
	}
}

func server(echoChan, senderChan chan string) {
	for {
		select {
		case data := <-echoChan:
			fmt.Printf("[RECV]: %v\n", data)
		case input := <-senderChan:
			fmt.Printf("[SENDER]: %v\n", input)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	recvBuf := make(chan string)
	inputBuf := make(chan string)

	conn, err := tcpListener()
	if err != nil {
		panic(err)
	}

	wg.Add(2)
	go tcpReceiver(conn, recvBuf, &wg)
	go tcpSender(conn, inputBuf, &wg)
	go server(recvBuf, inputBuf)
	wg.Wait()
}
