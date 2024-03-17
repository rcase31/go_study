package main

import (
	"fmt"
	"sync"
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

// behavior: same as above but using the sync package instead of channel
func WaitingForAllWithWaitGroup() {
	const NUM_WORKERS = 10
	var wg sync.WaitGroup

	for i := 0; i < NUM_WORKERS; i++ {
		// addition of workers will occur outside the go-routine because it needs to be synchronous
		wg.Add(1)
		go func() {
			// completion of the go-routine will be marked last and we use defer for it
			defer wg.Done()
			time.Sleep(time.Millisecond)
		}()
	}
	wg.Wait()
}

// It will successfully deplete the channel using a for loop (imperfect)
func WaitingForAllBufferedIterateChannel() {

	const NUM_WORKERS = 10

	done := make(chan bool, NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		go func() {
			worker(done)
		}()
	}
	// this is not ideal, because there's no guarantee that all workers will finish after this timing. Or it is rudimentary.
	time.Sleep(time.Second)
	// I do need to close the channel before iterating over it, otherwise it will never stop looping over it.
	close(done)
	for a := range done {
		fmt.Print(a)
	}
}

// behavior: execution will wait each go-routine to finish
func WaitingForAllBuffered() {

	const NUM_WORKERS = 100

	done := make(chan bool, NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		go worker(done)
	}
	// note that I do not NEED to close the channel here like in WaitingForAllBufferedIterateChannel
	for i := 0; i < NUM_WORKERS; i++ {
		<-done
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
