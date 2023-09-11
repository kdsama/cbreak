package cbreak

import "context"

type Action func() (interface{}, error)
type CircuitBreaker struct {
}

func New() *CircuitBreaker {
	return &CircuitBreaker{}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function

	res, err := fn()
	return res, err

}
