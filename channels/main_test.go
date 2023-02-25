package main

import (
	"testing"
)

func TestWaitingForAllUnbuffered(t *testing.T) {
	t.Run("WaitingForAllUnbuffered", func(t *testing.T) {
		WaitingForAllUnbuffered()
		t.Log("worked")
	},
	)
}

func TestWaitingForAllBufferedIterateChannel(t *testing.T) {
	t.Run("WaitingForAllBufferedIterateChannel", func(t *testing.T) {
		WaitingForAllBufferedIterateChannel()
		t.Log("worked")
	},
	)
}
