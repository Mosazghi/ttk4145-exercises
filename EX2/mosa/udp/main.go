package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const SENDER_PORT = 20023

const SERVER_IP = "10.100.23.11"
const SENDER_IP = SERVER_IP + ":200" + "23"
const ECHO_IP = SERVER_IP + ":200" + "23"
const BROADCAST_IP = "0.0.0.0" + ":30000"

func DialBroadcastUDP(port int) (net.PacketConn, error) {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		fmt.Println("Error: Socket:", err)
		return nil, err
	}
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		fmt.Println("Error: SetSockOpt REUSEADDR:", err)
		return nil, err
	}
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		fmt.Println("Error: SetSockOpt BROADCAST:", err)
		return nil, err
	}
	syscall.Bind(s, &syscall.SockaddrInet4{Port: port})
	if err != nil {
		fmt.Println("Error: Bind:", err)
		return nil, err
	}

	f := os.NewFile(uintptr(s), "")
	conn, err := net.FilePacketConn(f)
	if err != nil {
		fmt.Println("Error: FilePacketConn:", err)
		return nil, err
	}
	f.Close()

	return conn, nil
}

func udpSender(wg *sync.WaitGroup, inputChan chan string) {
	defer wg.Done()
	msgReader := bufio.NewReader(os.Stdin)
	conn, err := DialBroadcastUDP(20023)

	if err != nil {
		panic(err)
	}

	senderIpAddr, err := net.ResolveUDPAddr("udp", SENDER_IP)
	if err != nil {
		panic(err)
	}

	fmt.Println("[SENDER] Connected to ", conn.LocalAddr())

	defer conn.Close()
	for {
		// fmt.Print("> ")
		message, err := msgReader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		if err != nil {
			fmt.Println("Error reading message: ", err)
			continue
		}

		if len(message) == 0 {
			continue
		}

		_, err = conn.WriteTo([]byte(message), senderIpAddr)
		if err != nil {
			fmt.Println("ERROR: Write failed")
		}
		inputChan <- string(message)

		time.Sleep(10 * time.Millisecond)
	}
}

func udpBroadCastRecv(wg *sync.WaitGroup, broadCastChan chan<- string) {
	defer wg.Done()

	udpAddr, err := net.ResolveUDPAddr("udp", BROADCAST_IP)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("[BROADCAST] Connected to ", conn.LocalAddr())

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from udp: %v\n", err)
		}

		data := buffer[:n]
		broadCastChan <- string(data)
	}
}

func udpEchoRecv(wg *sync.WaitGroup, echoChan chan<- string) {
	defer wg.Done()

	conn, err := DialBroadcastUDP(20023)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Println("[ECHO] Connected to ", conn.LocalAddr())

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from udp: %v\n", err)
		}

		data := buffer[:n]
		echoChan <- string(data)
	}
}

func server(broadCastChan, senderChan, echoChan chan string) {
	for {
		select {
		case data := <-broadCastChan:
			fmt.Printf("[BROADCAST]: %v\n", data)
		case input := <-senderChan:
			fmt.Printf("[SENDER]: %v\n", input)
		case echo := <-echoChan:
			fmt.Printf("[ECHO]: %v\n", echo)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	broadCastChan := make(chan string)
	senderChan := make(chan string)
	echoChan := make(chan string)
	wg.Add(3)
	go udpSender(&wg, senderChan)
	go udpBroadCastRecv(&wg, broadCastChan)
	go udpEchoRecv(&wg, echoChan)
	go server(broadCastChan, senderChan, echoChan)
	wg.Wait()
}
