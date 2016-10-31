package main

import (
	"testing"
)

func TestStack(t *testing.T) {
	var stack StringStack

	if !stack.Empty() {
		t.Error("default value for stack is not empty")
	}

	stack.Push("test")

	if stack.Empty() {
		t.Error("non empty stack marked as empty")
	}

	if stack.Peek() != "test" {
		t.Error("incorrect peek")
	}

	if stack.Pop() != "test" {
		t.Error("incorrect pop")
	}

	if !stack.Empty() {
		t.Error("pop does not empty the stack")
	}
}
