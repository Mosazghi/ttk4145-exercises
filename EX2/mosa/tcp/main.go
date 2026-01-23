package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	ServerPort = "34933"
)

var localIP string

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

func tcpListener(serverIP string) (net.Conn, error) {
	fmt.Println("Listening...")
	addr, err := net.ResolveTCPAddr("tcp", serverIP+":"+ServerPort)
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

func tcpReceiver(ctx context.Context, cancel context.CancelFunc, conn net.Conn, recvChan chan string) {
	defer conn.Close()
	defer cancel()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF), errors.Is(err, io.ErrClosedPipe):
				fmt.Println("EOF")
			default:
				fmt.Println("Unknwon error, ", err.Error())
			}
			return 
		}
		data := string(buffer[:n])
		recvChan <- data

	}
}

func tcpSender(ctx context.Context, conn net.Conn, inputChan chan string) {
	msgReader := bufio.NewReader(os.Stdin)
	fmt.Println("[SENDER] Connected to ", conn.LocalAddr())
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {

		select {
		case <-ctx.Done():
			fmt.Println("Sender done")
			return
		default:
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
}


func main() {
	recvBuf := make(chan string)
	inputBuf := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ip, err := LocalIP()
	if err != nil {
		panic(err)
	}

	conn, err := tcpListener(ip)
	if err != nil {
		panic(err)
	}

	go tcpReceiver(ctx, cancel, conn, recvBuf)
	go tcpSender(ctx,  conn, inputBuf)
	for { 
		select {
			case <-ctx.Done():
				fmt.Println("Main done")
				return 
			case msg := <-recvBuf:
				fmt.Printf("[RECV] %s", msg)
			case msg := <-inputBuf:
				fmt.Printf("[SENT] %s", msg)

		}
	}
}
