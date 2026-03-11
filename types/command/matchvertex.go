package command

import (
	. "github.com/osingaatje/seshat/types/graph/parse-result"
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
	V1 *ParsedVertex
	V2 *ParsedVertex
	// can also put matching options in here such as matching with specific algorithms, should we need it
}
