package cbreak

type State byte

const (
	Open State = iota
	Half
	Closed
)
