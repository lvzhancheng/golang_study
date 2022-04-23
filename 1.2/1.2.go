package main

import (
	"fmt"
	"time"
)

func producer(data chan<- int) {
	defer close(data)
	for i := 1; i < 20; i++ {
		time.Sleep(time.Second)
		data <- i
	}
}
func consumer(data <-chan int, done chan bool) {
	for v := range data {
		// time.Sleep(time.Second)
		fmt.Println(v)
	}
	done <- true
}
func main() {
	data := make(chan int, 10)
	done := make(chan bool)
	go producer(data)
	// go producer(data)
	go consumer(data, done)
	<-done
}
