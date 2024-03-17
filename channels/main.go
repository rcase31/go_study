package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Millisecond)
	fmt.Println("done")
	done <- true
}

// behavior: successful execution
func OKBuffered() {

	done := make(chan bool, 1)
	done <- false
	<-done
	fmt.Println("execution success")
}

// behavior: deadlock, because execution will be halted when attempting to push value to channel
func DeadlockUnbuffered() {

	done := make(chan bool)
	done <- false
	<-done
	fmt.Println("never reaching this point")
}

// behavior: execution will wait each go-routine to finish
func WaitingForAllUnbuffered() {

	const NUM_WORKERS = 10

	done := make(chan bool)

	// go-calling all workers
	for i := 0; i < NUM_WORKERS; i++ {
		go worker(done)
	}
	for i := 0; i < NUM_WORKERS; i++ {
		<-done
	}
}

// behavior: execution will wait each go-routine to finish
func WaitingForAllBuffered() {

	const NUM_WORKERS = 100

	done := make(chan bool, NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		go worker(done)
	}
	for i := 0; i < NUM_WORKERS; i++ {
		<-done
	}
}

// It will successfully deplete the channel using a for loop
func WaitingForAllBufferedIterateChannel() {

	const NUM_WORKERS = 10

	done := make(chan bool, NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		go func() {
			worker(done)
		}()
	}
	time.Sleep(time.Second)

	close(done)
	for a := range done {
		fmt.Print(a)
	}
}

// behavior: execution will not wait each go-routine to finish
func ClosingTooEarly() {
	done := make(chan struct{}, 2)

	for i := 0; i < 2; i++ {
		go func() {
			time.Sleep(time.Second * 3)
			done <- struct{}{}
		}()
	}
	close(done)
	<-done
	<-done
}

// behavior: execution will panic because it was trying to send to a closed channel
func SendingToClosedChannel() {

	defer func() {
		didPanic := recover()
		if didPanic != nil {
			fmt.Println(didPanic)
		}
	}()

	done := make(chan bool, 2)
	close(done)
	done <- true
	<-done
}
