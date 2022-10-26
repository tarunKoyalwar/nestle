package nestle

type stack []int

// IsEmpty : check if stack is empty
func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto stack
func (s *stack) Push(i int) {
	*s = append(*s, i)
}

// Pop last element from stack and return it
func (s *stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return -1, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

// GetLastElement
func (s *stack) GetLastElement() (int, bool) {
	if s.IsEmpty() {
		return -1, false
	} else {
		return (*s)[len(*s)-1], true
	}
}
