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
		state:       Closed,
		threshold:   10,
		duration:    1 * time.Second,
		NotifyFunc:  n,
		hsThreshold: 5,
	}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function

	switch cb.state {
	case Closed:

		return cb.Run(fn)
	case Open:
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
	fmt.Println("Being called or not ?????")
	fmt.Println(cb.state)
	go cb.NotifyFunc(0)
	time.Sleep(cb.duration)
	cb.goodReqs = 0
	cb.state = Half
	go cb.NotifyFunc(1)

}

func (cb *CircuitBreaker) ReturnState() State {
	return cb.state
}
