package main

type input[T any] interface {
	Chan() chan T
	Error() error
}

type output[T any] interface {
	Write(t T)
	End(err error)
}

// chain implements input
var _ input[any] = newStream[any](0)

// chain implements output
var _ output[any] = newStream[any](0)

type chain[T any] struct {
	err error
	ch  chan T
}

// Error checks if there was an error in the chain
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

// End closes channel and sets error
func (s *chain[T]) End(err error) {
	s.err = err

	close(s.ch)
}
