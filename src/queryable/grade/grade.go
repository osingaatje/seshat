package grade

import (
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

func gradeDiag(c *context.Ctx, cmd command.GradeCmd) *grade.GradeResult {
	c.LogPrefixAdd("grade")
	defer c.LogPrefixRm("grade")

	possibleVertexMappings, certainties := getAlternativeSolutions(c, cmd)

	if len(possibleVertexMappings) == 0 {
		c.LogErr("No vertex mappings? Bug in code?")
	}

	// choose best alternative based on certainties:
	var bestScore float64 = -1
	var bestMapping map[VertexIdentifier]VertexIdentifier = nil

	for _, mp := range possibleVertexMappings {
		score := mappingScore(certainties, mp)
		if score > bestScore {
			bestScore = score
			bestMapping = mp
		}
	}

	if bestMapping == nil {
		panic("bug in code - fix your mapping choice")
	}

	// fix the vertices and grade the diagram
	// know nvertices are in the mapping (matched either by syntactic/semantic/structural matching)
	// extra vertices in the example solution means we have missing vertices in the submission
	// extra vertices in the submission should be treated as "erroneous" and points should be deducted
	// then, for all of the edges, do the same thing.

	c.LogErr("TODO MAKE GRADING WORK")
	return nil
}

func gradeVertex(cmd command.GradeCmd, refV VertexIdentifier, subV VertexIdentifier) {

}

func gradeEdge(cmd command.GradeCmd, refE EdgeIdentifier, subE EdgeIdentifier) {

}

func mappingScore(certainties map[VertexIdentifier]map[VertexIdentifier]float64, mapping map[VertexIdentifier]VertexIdentifier) float64 {
	var score float64 = 0
	for v1, v2 := range mapping {
		score += certainties[v1][v2]
	}
	return score
}
