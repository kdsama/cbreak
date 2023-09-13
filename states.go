package cbreak

// State type which tells which state the circuitbreaker is int
type State byte

const (
	//Open State
	Open State = iota
	//Half Open State
	Half
	//Closed State
	Closed
)

var (
	stateMapping = map[State]int{
		Open:   0,
		Half:   1,
		Closed: 2,
	}
)
