package main

type input[T any] interface {
	Chan() chan T
	Error() error
}

type output[T any] interface {
	Write(t T)
	End()
	Destroy(err error)
}

// chain implements input
var _ input[any] = newStream[any](0)

// chain implements output
var _ output[any] = newStream[any](0)

type chain[T any] struct {
	err error
	ch  chan T
}

// Error should be called after each receiving from the Chan()
func (s *chain[T]) Error() error {
	return s.err
}

func (s *chain[T]) Chan() chan T {
	return s.ch
}

func newStream[T any](cap int) *chain[T] {
	return &chain[T]{
		err: nil,
		ch:  make(chan T, cap),
	}
}

func (s *chain[T]) Write(t T) {
	s.ch <- t
}

func (s *chain[T]) End() {
	close(s.ch)
}

// Destroy should not be called after End() or vise-versa
func (s *chain[T]) Destroy(err error) {
	s.err = err

	close(s.ch)
}
