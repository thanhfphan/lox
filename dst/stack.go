package dst

type Stack[T any] struct {
	*Node[T]
}

type Node[T any] struct {
	Val  T
	Next *Node[T]
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Node: nil,
	}
}

func (s *Stack[T]) Len() int {
	if s.Node == nil {
		return 0
	}

	l := 0
	tmp := s.Node
	for tmp != nil {
		l++
		tmp = tmp.Next
	}

	return l
}

func (s *Stack[T]) Push(v T) {
	n := &Node[T]{
		Val:  v,
		Next: s.Node,
	}
	s.Node = n
}

func (s *Stack[T]) Pop() *Node[T] {
	if s.Node == nil {
		return nil
	}

	pop := s.Node
	s.Node = pop.Next
	return pop
}

func (s *Stack[T]) Peek() *Node[T] {
	if s.Node == nil {
		return nil
	}

	return s.Node
}

func (s *Stack[T]) IsEmpty() bool {
	return s.Node == nil
}
