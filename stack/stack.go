package stack

type Stack[T any] []T

func NewStack[T any]() Stack[T] {
	return Stack[T]{}
}

func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	last_i := len(*s) - 1
	val := (*s)[last_i]
	*s = (*s)[:last_i]
	return val, true
}
func (s *Stack[T]) Peek() (T, bool) {
	length := len(*s)

	if length == 0 {
		var zero T
		return zero, false
	}
	last_i := length - 1
	val := (*s)[last_i]
	return val, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Size() int {
	return len(*s)
}
func (s *Stack[T]) Contains(target T, equals func(a T, b T) bool) bool {
	for _, v := range *s {
		if equals(v, target) {
			return true
		}
	}
	return false
}

func (s *Stack[T]) Drain() []T {
	res := make([]T, 0, len(*s))
	for {
		if pop, ok := s.Pop(); !ok {
			break
		} else {
			res = append(res, pop)
		}
	}
	return res
}
