package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOKBuffered(t *testing.B) {
	assert.NotPanics(t, func() { OKBuffered() })
}
