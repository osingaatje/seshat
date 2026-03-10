package command

import (
	. "github.com/osingaatje/seshat/types/parse-result"
)

type MatchStringCmd struct {
	ref string // reference value
	act string // actual value
}

type MatchVertexCmd struct {
	v1 *ParsedVertex
	v2 *ParsedVertex
	// can also put matching options in here such as matching with specific algorithms, should we need it
}
