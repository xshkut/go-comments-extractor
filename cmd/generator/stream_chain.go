package main

type streamChain[T any] struct {
	err error
	ch  chan T
}

func NewStream[T any](cap int) *streamChain[T] {
	return &streamChain[T]{
		err: nil,
		ch:  make(chan T, cap),
	}
}

// Error should be called after each receiving from the Chan()
func (s *streamChain[T]) Error() error {
	return s.err
}

// Destroy should not be called after End() or vise-versa
func (s *streamChain[T]) Destroy(err error) {
	s.err = err

	close(s.ch)
}

func (s *streamChain[T]) End() {
	close(s.ch)
}

func (s *streamChain[T]) Chan() chan T {
	return s.ch
}

func (s *streamChain[T]) Write(t T) {
	s.ch <- t
}
