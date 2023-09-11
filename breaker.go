package cbreak

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit open")
)

type Action func() (interface{}, error)

type CircuitBreaker struct {
	count       int
	state       State
	hsThreshold int
	threshold   int
	goodReqs    int
	NotifyFunc  func(s int)
	duration    time.Duration
}

func New() *CircuitBreaker {
	return &CircuitBreaker{threshold: 10}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function
	switch cb.state {
	case Closed:
		cb.state = Open
		return nil, ErrCircuitOpen
	case Open:
		return nil, ErrCircuitOpen
	case Half:
		return cb.RunInHalfState(fn)
	}
	return nil, nil
}

func (cb *CircuitBreaker) RunInHalfState(fn Action) (interface{}, error) {

	res, err := fn()
	if err != nil {
		go cb.OpenCircuit()
		return res, err
	}
	cb.goodReqs++
	if cb.goodReqs >= cb.hsThreshold {
		cb.CloseCircuit()
	}
	return res, err
}
func (cb *CircuitBreaker) Run(fn Action) (interface{}, error) {

	res, err := fn()
	if err != nil {
		cb.count++
	}
	if cb.count >= cb.threshold {
		cb.OpenCircuit()
	}
	return res, err
}

func (cb *CircuitBreaker) CloseCircuit() {
	cb.state = Closed
	cb.goodReqs = 0
	go cb.NotifyFunc(0)

}

func (cb *CircuitBreaker) OpenCircuit() {
	cb.state = Open
	time.Sleep(cb.duration)
	cb.goodReqs = 0
	cb.state = Half
	go cb.NotifyFunc(1)

}
