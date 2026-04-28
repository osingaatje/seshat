package grade

import (
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

func gradeDiag(c *context.Ctx, cmd command.GradeCmd) *grade.GradeResult {
	c.LogPrefixAdd("grade '%s'", cmd.Submission.Metadata.Filename)
	defer c.LogPrefixRm("grade '%s'", cmd.Submission.Metadata.Filename)

	possibleVertexMaps, possibleEdgeMaps, certainties, err := getAlternativeSolutions(c, cmd)
	if err != nil {
		c.LogErr(err.Error())
		return nil
	}

	if len(possibleVertexMaps) == 0 {
		c.LogErr("No vertex mappings? Bug in code?")
		return nil
	}
	if len(possibleVertexMaps) != len(possibleEdgeMaps) {
		c.LogErr("BUG IN CODE - Vertex Maps and Edge Maps not equal lengths")
		return nil
	}

	// choose best alternative based on certainties:
	var bestScore float64 = -1
	var finalMappingIndex int = -1

	for i := range len(possibleVertexMaps) {
		score := mappingScore(certainties, possibleVertexMaps[i], possibleEdgeMaps[i])
		if score > bestScore {
			bestScore = score
			finalMappingIndex = i
		}
	}

	if finalMappingIndex == -1 {
		panic("bug in code - fix your mapping choice")
	}

	// finalVertexMapping := possibleVertexMaps[finalMappingIndex]
	// finalEdgeMapping := possibleEdgeMaps[finalMappingIndex]

	// fix the vertices and grade the diagram
	// known vertices are in the mapping (matched either by syntactic/semantic/structural matching)
	// extra vertices in the example solution means we have missing vertices in the submission
	// extra vertices in the submission should be treated as "erroneous" and points should be deducted
	// then, for all of the edges, do the same thing.
	// result := grade.NewGradeResult()
	// for refEId, refE := range cmd.ReferenceSolution.Edges {
	// }

	c.LogErr("TODO MAKE GRADING WORK")
	return nil
}

func gradeVertex(cmd command.GradeCmd, refV VertexIdentifier, subV VertexIdentifier) {

}

func gradeEdge(cmd command.GradeCmd, refE EdgeIdentifier, subE EdgeIdentifier) {

}

func mappingScore(
	certainties map[VertexIdentifier]map[VertexIdentifier]float64,
	vertMap map[VertexIdentifier]VertexIdentifier,
	edgeMap map[EdgeIdentifier]EdgeIdentifier,
) float64 {
	var score float64 = 0
	for v1, v2 := range vertMap {
		score += certainties[v1][v2]
	}

	score += float64(len(edgeMap)) / 2 // add .5 for every discovered edge

	return score
}
