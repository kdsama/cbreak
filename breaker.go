package cbreak

import (
	"context"
	"errors"
	"fmt"
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

func New(n func(s int)) *CircuitBreaker {

	return &CircuitBreaker{
		state:      Closed,
		threshold:  10,
		duration:   5 * time.Second,
		NotifyFunc: n,
	}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function

	switch cb.state {
	case Closed:

		return cb.Run(fn)
	case Open:
		fmt.Println("State", cb.state)
		return nil, ErrCircuitOpen
	case Half:
		return cb.RunInHalfState(fn)
	}
	return cb.Run(fn)
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
		go cb.OpenCircuit()
	}
	return res, err
}

func (cb *CircuitBreaker) CloseCircuit() {
	cb.state = Closed
	cb.goodReqs = 0
	go cb.NotifyFunc(2)

}

func (cb *CircuitBreaker) OpenCircuit() {
	cb.state = Open
	go cb.NotifyFunc(2)
	time.Sleep(cb.duration)
	cb.goodReqs = 0
	cb.state = Half
	go cb.NotifyFunc(1)

}

func (cb *CircuitBreaker) ReturnState() int {
	switch cb.state {
	case Open:
		return 2
	case Half:
		return 1
	case Closed:
		return 0
	}
	return -1
}
