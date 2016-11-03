package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)
	var stack StringStack

	assert.True(stack.Empty(), "default value for a stack should be empty")
	assert.Equal("", stack.Pop(), "pop value for empty stack should be empty string")
	assert.Equal("", stack.Peek(), "peek value for empty stack should be empty string")

	stack.Push("foo")

	assert.False(stack.Empty(), "after a push the stack should not be empty")
	assert.True(stack.Contains("foo"), "stack should contain value after push")
	assert.False(stack.Contains("bar"), "stack should not contain values that were not pushed")
	assert.Equal("foo", stack.Peek(), "peek value should equal the latest push")
	assert.Equal("foo", stack.Pop(), "pop value should equal the latest push")
	assert.True(stack.Empty(), "after pop the stack should be empty")

	for i := 0; i < 10; i++ {
		stack.Push(" ")
	}

	assert.Equal(10, stack.Size(), "stack should contain N elements after N pushes")
}
