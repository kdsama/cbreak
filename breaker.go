package cbreak

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	//ErrCircuitOpen is returned on function calls during Open Circuit State
	ErrCircuitOpen = errors.New("circuit open")
)

// Action is the signature used inside Execute function that runs the circuitBreaker.
type Action func() (interface{}, error)

// CircuitBreakerOpts is the option Configuration struct
type CircuitBreakerOpts struct {
	// Number of errors past which the circuit breaker will open
	Threshold uint

	//
	// value from 1 to 99. When HsThresholdPercentage of good requests is reached in half state , circuit is closed
	// any input below 1 will be set to 1 and value above 99 will be set to 99
	HsThresholdPercentage uint

	// Function passed with this signature will be called whenever a state is changed.
	NotifyFunc func(s int)

	//Sleep time after which OpenCircuit switches to Half-Open.
	Duration time.Duration
}

// CircuitBreaker is the main struct of the project
type CircuitBreaker struct {
	count       int
	hsThreshold int
	threshold   int
	goodReqs    int

	// lock for state and good requests
	sLock *sync.Mutex
	state State

	notifyFunc func(s int)
	r          uint32
	duration   time.Duration
}

// New initiate a new CircuitBreaker Object.
func New(opts CircuitBreakerOpts) *CircuitBreaker {
	if opts.Threshold < 1 {
		opts.Threshold = 1
	}
	if opts.HsThresholdPercentage > 99 {
		opts.HsThresholdPercentage = 99
	}
	if opts.HsThresholdPercentage < 1 {
		opts.HsThresholdPercentage = 1
	}
	return &CircuitBreaker{
		state:       Closed,
		sLock:       &sync.Mutex{},
		threshold:   int(opts.Threshold),
		duration:    opts.Duration,
		notifyFunc:  opts.NotifyFunc,
		hsThreshold: int(opts.HsThresholdPercentage),
	}
}

// Execute executes the user defined function in the circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, fn Action) (interface{}, error) {
	// Execute the function
	var state State

	cb.sLock.Lock()
	state = cb.state
	cb.sLock.Unlock()

	switch state {
	case Closed:

		return cb.run(fn)
	case Open:
		return nil, ErrCircuitOpen
	case Half:
		return cb.runInHalfState(fn)
	}
	return cb.run(fn)
}

// GetState returns current state of the circuit breaker. 0 -> Open, 1 -> Half, 2 -> Closed
func (cb *CircuitBreaker) GetState() State {
	var state State

	cb.sLock.Lock()
	state = cb.state
	cb.sLock.Unlock()

	return state
}

func (cb *CircuitBreaker) runInHalfState(fn Action) (interface{}, error) {

	res, err := fn()
	if err != nil {
		cb.openCircuit()
		return res, err
	}

	cb.goodReqs++

	if cb.goodReqs >= (cb.hsThreshold*cb.threshold)/100 {
		cb.closeCircuit()
	}
	return res, err
}
func (cb *CircuitBreaker) run(fn Action) (interface{}, error) {

	res, err := fn()
	if err != nil {
		cb.count++
	}
	if cb.count >= cb.threshold {
		go cb.openCircuit()
	}
	return res, err
}

func (cb *CircuitBreaker) closeCircuit() {
	cb.setState(Closed)

}

func (cb *CircuitBreaker) openCircuit() {
	cb.setState(Open)
	go cb.halfCircuit()

}
func (cb *CircuitBreaker) halfCircuit() {
	// If value of r is not same after the sleep, means a new goroutine has been invoked.
	//Hence we dont need to set state anymore
	atomic.AddUint32(&cb.r, 1)
	a := cb.r
	time.Sleep(cb.duration)
	if a != cb.r {
		return
	}

	cb.setState(Half)
}

func (cb *CircuitBreaker) setState(st State) {

	cb.sLock.Lock()
	cb.goodReqs = 0
	cb.count = 0
	cb.state = st
	cb.sLock.Unlock()

	//Notify the userDefined Function
	go cb.notifyFunc(stateMapping[st])

}
