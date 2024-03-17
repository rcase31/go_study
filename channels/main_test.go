package main

import (
	"testing"
	"time"
)

func TestWaitingForAllUnbuffered(t *testing.T) {
	t.Run("Waiting for all in an unbuffered channel", func(t *testing.T) {
		WaitingForAllUnbuffered()
		t.Log("worked")
	},
	)
}

func TestWaitingForAllBufferedIterateChannel(t *testing.T) {
	t.Run("Waiting for all buffered: iterate over channel", func(t *testing.T) {
		WaitingForAllBufferedIterateChannel()
		t.Log("worked")
	},
	)
}

func TestDeadLockUnbuffered(t *testing.T) {
	t.Run("Deadlock Unbuffered Channel", func(t *testing.T) {
		chFinished := make(chan struct{})
		go func() {
			DeadlockUnbuffered()
			chFinished <- struct{}{}
		}()
		select {
		case <-chFinished:
			t.Fatal("test should not have finished successuly")
		case <-time.After(time.Second):
			t.Log("test cause an expected deadlock")
		}

	})
}
