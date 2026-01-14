package main

import (
	"fmt"
	"time"
)

func producer(bufferChan chan<- int) {
	defer close(bufferChan)
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[producer]: pushing %d\n", i)
		bufferChan <- i
	}
}

func consumer(bufferChan <-chan int, done chan<- bool) {
	time.Sleep(1 * time.Second)
	for i := range bufferChan {
		fmt.Printf("[consumer]: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}

	done <- true
}

func main() {
	bufferChan := make(chan int, 5)
	done := make(chan bool)
	go consumer(bufferChan, done)
	go producer(bufferChan)

	<-done
}
