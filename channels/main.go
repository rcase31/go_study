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
	done := make(chan bool, 2)

	for i := 0; i < 2; i++ {
		go func() {
			time.Sleep(time.Second)
			done <- true
		}()
	}
	close(done)
	fmt.Println(<-done)
	fmt.Println(<-done)
	time.Sleep(time.Second)
}

// behavior: execution will panic because it was trying to send to a closed channel
func SendingToClosedChannel() {
	done := make(chan bool, 2)
	close(done)
	done <- true
	<-done
}

// behavior: it will print default boolean value (false)
func WorkingWithEmptyChannel() {
	done := make(chan bool, 2)

	close(done)
	fmt.Println(len(done))
	fmt.Println(<-done)
	fmt.Println(<-done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}

//TODO: make test where I can test if a go routine still excutes after function is done
// However, if the main function of your program exits, or if os.Exit() is called, the program will terminate, and all running goroutines will be stopped abruptly, regardless of their state. This means that if your main program's execution completes while your goroutines are still working, those goroutines will be stopped, and they won't complete their execution.

//TODO: make a case where I use defer to guarantee that all is done

//TODO: use a channel as an atomic counter( using len)

//TODO: make another package to test atomicity
