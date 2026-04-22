package command

import (
	. "github.com/osingaatje/seshat/types/graph/intern"
)

type MatchStringCmd struct {
	Ref string // reference value
	Act string // actual value
}

type MatchStringRes struct {
	Score float64
	Err   error
}

type MatchVertexCmd struct {
	V1 *InternalVertex
	V2 *InternalVertex
	// can also put matching options in here such as matching with specific algorithms, should we need it
}
