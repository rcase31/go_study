package main

import (
	"testing"
	"time"
)

func shouldTimeout(t *testing.T, fn func()) {
	t.Helper()
	chFinished := make(chan struct{})
	go func() {
		fn()
		chFinished <- struct{}{}
	}()
	select {
	case <-chFinished:
		t.Fatal("test should not have finished successuly")
	case <-time.After(time.Second * 2):
		t.Log("test cause an expected deadlock")
	}
}

func shouldNotTimeout(t *testing.T, fn func()) {
	t.Helper()
	chFinished := make(chan struct{})
	go func() {
		fn()
		chFinished <- struct{}{}
	}()
	select {
	case <-chFinished:
		t.Log("test should have finished successuly")
	case <-time.After(time.Second * 2):
		t.Fatal("test should not time-out")
	}
}

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
	},
	)
}

func TestWaitAllWorkersWithoutChannel(t *testing.T) {
	t.Run("We wait for all go routines without using any channels", func(t *testing.T) {
		shouldNotTimeout(t, WaitingForAllWithWaitGroup)
	})
}

func TestDeadLockUnbuffered(t *testing.T) {
	t.Run("Deadlock Unbuffered Channel", func(t *testing.T) {
		shouldTimeout(t, DeadlockUnbuffered)
	})
}
