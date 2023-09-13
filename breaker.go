package cbreak

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit open")
)

type Action func() (interface{}, error)

type CircuitBreaker struct {
	count int

	sLock       *sync.Mutex
	state       State
	hsThreshold int
	threshold   int
	goodReqs    int
	NotifyFunc  func(s int)
	r           uint32
	duration    time.Duration
}

func New(n func(s int)) *CircuitBreaker {

	return &CircuitBreaker{
		state:       Closed,
		sLock:       &sync.Mutex{},
		threshold:   10,
		duration:    1 * time.Second,
		NotifyFunc:  n,
		hsThreshold: 5,
	}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function
	var state State
	cb.sLock.Lock()

	state = cb.state

	cb.sLock.Unlock()
	switch state {
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
		cb.OpenCircuit()
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
	cb.SetState(Closed)

}

func (cb *CircuitBreaker) OpenCircuit() {
	cb.SetState(Open)
	go cb.HalfCircuit()

}
func (cb *CircuitBreaker) HalfCircuit() {
	atomic.AddUint32(&cb.r, 1)
	a := cb.r
	time.Sleep(cb.duration)
	if a != cb.r {
		return
	}
	cb.SetState(Half)
}

func (cb *CircuitBreaker) SetState(st State) {

	cb.sLock.Lock()
	cb.goodReqs = 0
	cb.state = st
	cb.sLock.Unlock()
	go cb.NotifyFunc(stateMapping[st])

}

func (cb *CircuitBreaker) ReturnState() State {
	var state State
	cb.sLock.Lock()
	state = cb.state
	cb.sLock.Unlock()
	return state
}
