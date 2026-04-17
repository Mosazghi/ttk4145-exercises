package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

const (
	BROADCAST_IP = "localhost:30000"
	SENDER_PORT  = 30000
)

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

func backup(switchStateSignal chan int) {
	count := 0
	fmt.Print("Backup start")

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
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				fmt.Printf("Read timeout: %v. Closing connection.\n", err)
				break
			}
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		data := buffer[:n]
		count, _ = strconv.Atoi(string(data))
	}

	switchStateSignal <- count
	fmt.Println("[Backup]: Switching to primary!")
}

func primary(initVal int) {
	fmt.Println("Primary start")
	count := initVal
	cmd := exec.Command("kitty", "go", "run", "pb.go")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("Started backup!")
	conn, err := DialBroadcastUDP(30001)
	if err != nil {
		panic(err)
	}

	senderIpAddr, err := net.ResolveUDPAddr("udp", BROADCAST_IP)
	if err != nil {
		panic(err)
	}

	fmt.Println("[SENDER] Connected to ", conn.LocalAddr())

	defer conn.Close()
	for {
		count++
		_, err = conn.WriteTo([]byte(strconv.Itoa(count)), senderIpAddr)
		if err != nil {
			fmt.Println("ERROR: Write failed")
		}

		fmt.Println("current count: ", count)
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	switchStateSignal := make(chan int)

	go backup(switchStateSignal)

	for {
		select {
		case singalCount := <-switchStateSignal:
			fmt.Println("[Server]: Switching to primary mode")
			go primary(singalCount)
		}
	}
}
