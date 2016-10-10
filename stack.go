package main

type StringStack struct {
	stack []string
	i     int
}

func (s *StringStack) Push(str string) {
	if s.i < len(s.stack) {
		s.stack[s.i] = str
	} else {
		s.stack = append(s.stack, str)
	}
	s.i++
}

func (s *StringStack) Pop() string {
	if s.i > 0 {
		s.i--
		return s.stack[s.i]
	}

	return ""
}

func (s *StringStack) Peek() string {
	if s.i > 0 {
		return s.stack[s.i-1]
	}

	return ""
}

func (s *StringStack) Empty() bool {
	return s.Size() <= 0
}

func (s *StringStack) Size() int {
	return s.i
}
