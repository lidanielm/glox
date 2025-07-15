package tool

type Stack[T any] struct {
	stack []T
	length int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{stack: []T{}, length: 0}
}

func (s *Stack[T]) Push(element T) {
	s.stack = append(s.stack, element)
	s.length++
}

func (s *Stack[T]) Pop() T {
	element := s.stack[s.length-1]
	s.stack = s.stack[:s.length-1]
	s.length--
	return element
}

func (s *Stack[T]) Peek() T {
	return s.stack[s.length - 1]
}

func (s *Stack[T]) Get(index int) T {
	return s.stack[index]
}

func (s *Stack[T]) Length() int {
	return s.length
}

func (s *Stack[t]) IsEmpty() bool {
	return s.length == 0
}