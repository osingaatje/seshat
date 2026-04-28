package grade

import (
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

func gradeDiag(c *context.Ctx, cmd command.GradeCmd) *GradeCalculation {
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

	// choose best vertex+edge mapping based on certainties:
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

	finalVertexMapping := possibleVertexMaps[finalMappingIndex]
	finalEdgeMapping := possibleEdgeMaps[finalMappingIndex]

	// fix the vertices and grade the diagram
	// known vertices are in the mapping (matched either by syntactic/semantic/structural matching)
	// extra vertices in the example solution means we have missing vertices in the submission
	// extra vertices in the submission should be treated as "erroneous" and points should be deducted
	// then, for all of the edges, do the same thing

	gradedSubVertices := map[VertexIdentifier]bool{}
	gradedSubEdges := map[EdgeIdentifier]bool{}

	result := NewGradeCalculation()
	for refEId, refE := range cmd.ReferenceSolution.Edges {
		subEId, ok := finalEdgeMapping[refEId]
		if !ok {
			// missing edge
			// TODO deduct points
			continue
		}
		gradeEdge(cmd, &result, refE, cmd.Submission.Edges[subEId])
		gradedSubEdges[subEId] = true
	}

	for refVId, refV := range cmd.ReferenceSolution.Vertices {
		subVId, ok := finalVertexMapping[refVId]
		if !ok {
			// missing vertex
			// TODO deduct points
			continue
		}
		gradeVertex(cmd, &result, refV, cmd.Submission.Vertices[subVId])
		gradedSubVertices[subVId] = true
	}

	for subEId, _ := range cmd.Submission.Edges {
		if _, ok := gradedSubEdges[subEId]; ok {
			continue
		}
		// extra edge that should have points deducted
		// TODO deduct points
	}
	for subVId, _ := range cmd.Submission.Vertices {
		if _, ok := gradedSubVertices[subVId]; ok {
			continue
		}
		// extra vertex that should have points deducted
		// TODO deduct points
	}

	result.EdgeGrades = nil // todo sdflksjdflkdsjfljdkf make grading

	return nil // &result
}

func gradeVertex(cmd command.GradeCmd, res *GradeCalculation, refV *InternalVertex, subV *InternalVertex) {
	if res == nil || refV == nil || subV == nil {
		panic("BUG CODE - gradeVertex nil vertex")
	}

	typeScore := Grade{}

	titleScore := Grade{}

	attrScores := map[string]AttributeScore{}
	// TODO traverse attributes

	res.VertexGrades[subV.Id] = GradeResultVertex{
		PresenceScore: Grade{
			Grade:  cmd.Rubric.Scores.VertexScore.PointsForPresence,
			Reason: GRADE_REASON_PRESENT,
		},
		TypeScore:       typeScore,
		TitleScore:      titleScore,
		AttributeScores: attrScores,
	}
}

func gradeEdge(cmd command.GradeCmd, res *GradeCalculation, refE *InternalEdge, subE *InternalEdge) {
	if res == nil || refE == nil || subE == nil {
		panic("BUG CODE - gradeEdge nil edge")
	}
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
