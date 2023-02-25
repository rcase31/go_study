package main

import (
	"fmt"
	"log"
	"time"
)

// TODO: use flag package here
func main() {
	start := time.Now()

	OKBuffered()
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
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

// TODO: create mechanics to ensure that we're done with this channel
func WaitingForAllBufferedIterateChannel() {

	const NUM_WORKERS = 10

	done := make(chan bool, NUM_WORKERS)
	counter := 0

	for i := 0; i < NUM_WORKERS; i++ {
		go func() {
			worker(done)
			//TODO: fix this because it is susceptivle to race conditions
			counter++
		}()
	}

	//TODO:
	for counter < NUM_WORKERS {
		time.Sleep(time.Second)
	}

	close(done)
	for a := range done {
		fmt.Print(a)
	}
}

// behavior: execution will not wait each go-routine to finish
func closingTooEarly() {
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
