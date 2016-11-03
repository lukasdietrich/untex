package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)
	var stack StringStack

	assert.True(stack.Empty(), "default value for a stack should be empty")

	stack.Push("test")

	assert.False(stack.Empty(), "after a push the stack should not be empty")
	assert.Equal("test", stack.Peek(), "peek value should equal the latest push")
	assert.Equal("test", stack.Pop(), "pop value should equal the latest push")
	assert.True(stack.Empty(), "after pop the stack should be empty")

	for i := 0; i < 10; i++ {
		stack.Push(" ")
	}

	assert.Equal(10, stack.Size(), "stack should contain N elements after N pushes")
}
