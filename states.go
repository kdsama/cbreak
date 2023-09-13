package cbreak

type State byte

const (
	Open State = iota
	Half
	Closed
)

var (
	stateMapping = map[State]int{
		Open:   0,
		Half:   1,
		Closed: 2,
	}
)
