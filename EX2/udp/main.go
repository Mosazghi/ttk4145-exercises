package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func udpSender(wg *sync.WaitGroup, inputChan chan string) {
	defer wg.Done()
	msgReader := bufio.NewReader(os.Stdin)
	udpAddr, err := net.ResolveUDPAddr("udp", "10.22.124.123:20000")
	if err != nil {
		panic(err)
	}

	// Dial the connection (binds a local port and sets the remote address)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to ", conn.LocalAddr())

	defer conn.Close()
	for {
		fmt.Print("> ")
		message, err := msgReader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		if err != nil {
			fmt.Println("Error reading message: ", err)
			continue
		}

		if len(message) == 0 {
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("ERROR: Write failed")
		}
		inputChan <- string(message)

		time.Sleep(10 * time.Millisecond)
	}
}

func udpEchoRecv(wg *sync.WaitGroup, bufChan chan<- string) {
	defer wg.Done()

	udpAddr, err := net.ResolveUDPAddr("udp", "10.22.124.123:30000")
	if err != nil {
		panic(err)
	}

	// Dial the connection (binds a local port and sets the remote address)
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("[RECV] Connected to ", conn.LocalAddr())

	buffer := make([]byte, 1024)

	for {
		fmt.Println("Reading...")
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from udp: %v\n", err)
		}
		fmt.Println("Read")

		data := buffer[:n]
		bufChan <- string(data)
		time.Sleep(500 * time.Millisecond)
	}
}

func server(bufChan, inputChan chan string) {
	for {
		select {
		case data := <-bufChan:
			fmt.Printf("Recv: %v\n", data)
		case input := <-inputChan:
			fmt.Printf("Sending: %v\n", input)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	bufChan := make(chan string)
	inputChan := make(chan string)
	// go udpSender()
	wg.Add(2)
	go udpSender(&wg, inputChan)
	go udpEchoRecv(&wg, bufChan)
	go server(bufChan, inputChan)
	wg.Wait()
}
