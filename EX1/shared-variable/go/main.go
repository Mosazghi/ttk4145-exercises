package main

import (
	. "fmt"
)

type Request int

var i = 0

const (
	Inc Request = iota
	Dec
)

var requestName = map[Request]string{
	Inc: "Increment",
	Dec: "Decrement",
}

func incrementing(c chan<- Request, done chan<- bool) {
	for range 1_000_000 {
		c <- Inc
	}
	done <- true
}

func decrementing(c chan<- Request, done chan<- bool) {
	for range 1_000_000 {
		c <- Dec
	}
	done <- true
}

func main() {
	channel := make(chan Request)
	done := make(chan bool)
	get := make(chan int)

	go incrementing(channel, done)
	go decrementing(channel, done)

	go func() {
		for {
			select {
			case msg := <-channel:
				switch msg {
				case Inc:
					i++
				case Dec:
					i--
				}
			case <-get:
				return
			}
		}
	}()

	<-done
	<-done

	get <- 1
	Println("The magic number is:", i)
}
