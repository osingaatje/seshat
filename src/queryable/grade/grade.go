package grade

import (
	"fmt"

	"github.com/osingaatje/seshat/helper/multiplicity"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

func gradeDiag(c *context.Ctx, cmd command.GradeCmd) *GradeCalculation {
	if cmd.Rubric == nil || cmd.ReferenceSolution == nil || cmd.Submission == nil {
		forFileName := ""
		if cmd.ReferenceSolution != nil {
			forFileName = " for solution: " + cmd.ReferenceSolution.Metadata.Filename
		}
		c.LogErr("Grader config, rubric, reference solution or submission was missing%s, returning nothing...", forFileName)
		return nil
	}

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
			result.MissingReferenceEdges[refEId] = GradeResultEdgeMissing(cmd.Rubric.Scores.EdgeScore.AbsenceIncompleteScore)
			continue
		}
		gradeEdge(c, cmd, &result, refE, cmd.Submission.Edges[subEId])
		gradedSubEdges[subEId] = true
	}

	for refVId, refV := range cmd.ReferenceSolution.Vertices {
		subVId, ok := finalVertexMapping[refVId]
		if !ok && cmd.Rubric.Scores.VertexScore.AbsenceIncompleteScore != 0 {
			// missing vertex
			result.MissingReferenceVertices[refVId] = GradeResultVertexMissing(cmd.Rubric.Scores.VertexScore.AbsenceIncompleteScore)
			continue
		}
		gradeVertex(c, cmd, &result, refV, cmd.Submission.Vertices[subVId])
		gradedSubVertices[subVId] = true
	}

	for subEId, _ := range cmd.Submission.Edges {
		if _, ok := gradedSubEdges[subEId]; ok {
			continue
		}
		// skip non-impacting scores
		if cmd.Rubric.Scores.VertexScore.SuperfluousScore == 0 {
			continue
		}
		// extra edge that should have points deducted
		result.EdgeGrades[subEId] = GradeResultEdgeExtra(cmd.Rubric.Scores.EdgeScore.SuperfluousScore)
	}

	for subVId, _ := range cmd.Submission.Vertices {
		if _, ok := gradedSubVertices[subVId]; ok {
			continue
		}
		if cmd.Rubric.Scores.VertexScore.SuperfluousScore == 0 {
			continue
		}
		// extra vertex that should have points deducted
		result.VertexGrades[subVId] = GradeResultVertexExtra(cmd.Rubric.Scores.VertexScore.SuperfluousScore)
	}

	return &result
}

func gradeVertex(c *context.Ctx, cmd command.GradeCmd, res *GradeCalculation, refV *InternalVertex, subV *InternalVertex) {
	if res == nil || refV == nil || subV == nil {
		panic("BUG CODE - gradeVertex nil vertex")
	}

	typeScore := calculateVertexTypeScore(c, cmd, refV, subV)
	titleScore := calculateVertexTitleScore(c, cmd, refV, subV)
	attrScores := calcalateVertexAttrScores(c, cmd, refV, subV)

	res.VertexGrades[subV.Id] = GradeResultVertex{
		PresenceScore: Grade{
			Grade:  cmd.Rubric.Scores.VertexScore.PresenceScore,
			Reason: GRADE_REASON_PRESENT,
		},
		TypeScore:       &typeScore,
		TitleScore:      &titleScore,
		AttributeScores: attrScores,
	}
}

func calculateVertexTypeScore(c *context.Ctx, cmd command.GradeCmd, refV *InternalVertex, subV *InternalVertex) Grade {
	return syntacticSemanticMatch(c,
		refV.Properties.Type, subV.Properties.Type, cmd.Rubric.GraderConfig.ClassContentSimilarity,
		cmd.Rubric.Scores.VertexTypeScore,
	)
}

func calculateVertexTitleScore(c *context.Ctx, cmd command.GradeCmd, refV *InternalVertex, subV *InternalVertex) Grade {
	return syntacticSemanticMatch(c,
		refV.Title, subV.Title, cmd.Rubric.GraderConfig.ClassContentSimilarity,
		cmd.Rubric.Scores.VertexTitleScore,
	)
}

func calcalateVertexAttrScores(c *context.Ctx, cmd command.GradeCmd, refV *InternalVertex, subV *InternalVertex) map[string]GradeResultVertexAttr {
	res := map[string]GradeResultVertexAttr{}

	for refName, refVal := range refV.Values {
		subVal, ok := subV.Values[refName]
		if !ok { // possible TODO: make this not add a grade if the penalty is 0
			res[refName] = GradeResultAttrMissing(cmd.Rubric.Scores.VertexAttributeScore.AbsenceIncompleteScore)
		}

		attrTypeScore := calculateVertexAttrTypeScore(c, cmd, refVal.Properties.Type, subVal.Properties.Type)
		attrVisibilityScore := calculateVertexAttrVisScore(c, cmd, refVal.Properties.Visibility, subVal.Properties.Visibility)

		res[refName] = GradeResultVertexAttr{
			NameScore:       Grade{Grade: cmd.Rubric.Scores.VertexAttributeScore.PresenceScore, Reason: GRADE_REASON_PRESENT},
			TypeScore:       &attrTypeScore,
			VisibilityScore: &attrVisibilityScore,
		}
	}

	for subName, _ := range subV.Values {
		if _, ok := res[subName]; ok {
			continue
		}
		if cmd.Rubric.Scores.VertexAttributeScore.SuperfluousScore == 0 {
			continue
		}
		// superfluous attribute:
		res[subName] = GradeResultAttrExtra(cmd.Rubric.Scores.VertexAttributeScore.SuperfluousScore)
	}

	return res
}

func calculateVertexAttrTypeScore(c *context.Ctx, cmd command.GradeCmd, refType string, subType string) Grade {
	return syntacticSemanticMatch(c, refType, subType, cmd.Rubric.GraderConfig.ClassContentSimilarity, cmd.Rubric.Scores.VertexTypeScore)
}
func calculateVertexAttrVisScore(c *context.Ctx, cmd command.GradeCmd, refVisibility ValuePropVisibilityVar, subVisibility ValuePropVisibilityVar) Grade {
	return formatPresentOrAbsenceGrade(refVisibility == subVisibility, cmd.Rubric.Scores.VertexTypeScore, GradeReasonIncorrect(string(refVisibility)))
}

func gradeEdge(c *context.Ctx, cmd command.GradeCmd, res *GradeCalculation, refE *InternalEdge, subE *InternalEdge) {
	if res == nil || refE == nil || subE == nil {
		panic("BUG CODE - gradeEdge nil edge")
	}

	edgeStyleScore := calculateEdgeTypeScore(cmd, refE, subE)
	edgeLabelScore := calculateEdgeLabelScore(c, cmd, refE.Label, subE.Label)
	edgeFromScore := calculateEdgeEndScore(c, cmd, refE, subE, refE.FromProperties, subE.FromProperties)
	edgeToScore := calculateEdgeEndScore(c, cmd, refE, subE, refE.ToProperties, subE.ToProperties)

	res.EdgeGrades[subE.Id] = GradeResultEdge{
		PresenceScore:  Grade{Grade: cmd.Rubric.Scores.EdgeScore.PresenceScore, Reason: GRADE_REASON_PRESENT},
		EdgeTypeScore:  &edgeStyleScore,
		EdgeLabelScore: &edgeLabelScore,
		StartScore:     &edgeFromScore,
		EndScore:       &edgeToScore,
	}
}

func calculateEdgeTypeScore(cmd command.GradeCmd, ref *InternalEdge, act *InternalEdge) Grade {
	return formatPresentOrAbsenceGradeWithCorrectMsg(
		ref.FromProperties.ArrowStyle == act.FromProperties.ArrowStyle &&
			ref.StyleProperties.LineStyle == act.StyleProperties.LineStyle &&
			ref.ToProperties.ArrowStyle == act.ToProperties.ArrowStyle,
		cmd.Rubric.Scores.EdgeTypeScore,
		GRADE_REASON_EQUAL,
		GRADE_REASON_INCORRECT,
	)
}
func calculateEdgeLabelScore(c *context.Ctx, cmd command.GradeCmd, refLbl *Label, actLbl *Label) Grade {
	if refLbl == nil {
		return formatPresentOrAbsenceGradeWithCorrectMsg(
			actLbl == nil,
			cmd.Rubric.Scores.EdgeLabelScore,
			GRADE_REASON_EQUAL,
			GRADE_REASON_INCORRECT,
		)
	}
	if actLbl == nil {
		return Grade{
			Grade:  cmd.Rubric.Scores.EdgeLabelScore.AbsenceIncompleteScore,
			Reason: GRADE_REASON_ABSENT,
		}
	}

	return syntacticSemanticMatch(c, refLbl.Text, actLbl.Text, cmd.Rubric.GraderConfig.ClassContentSimilarity, cmd.Rubric.Scores.EdgeLabelScore)
}

func calculateEdgeEndScore(c *context.Ctx, cmd command.GradeCmd, refEdge *InternalEdge, subEdge *InternalEdge, ref EdgeEndProperties, act EdgeEndProperties) Grade {
	score := cmd.Rubric.Scores.EdgeEndLabelScore
	if ref.Label == nil {
		if act.Label == nil {
			return Grade{Grade: score.PresenceScore, Reason: GRADE_REASON_EQUAL}
		}
		return Grade{Grade: score.SuperfluousScore, Reason: GRADE_REASON_SUPERFLUOUS}
	}

	if act.Label == nil {
		return Grade{Grade: score.AbsenceIncompleteScore, Reason: GRADE_REASON_ABSENT}
	}

	// figure out if we need to grade a multiplicity or something else:
	refMult, ok := multiplicity.GetMultiplicity(ref.Label.Text)
	if ok {
		actMult, ok := multiplicity.GetMultiplicity(act.Label.Text)
		if !ok {
			return Grade{Grade: score.AbsenceIncompleteScore, Reason: GradeReasonIncorrectMultiplicity("", ref.Label.Text)}
		}

		correctnessPercentage := refMult.Equal(actMult)
		if correctnessPercentage < 0.01 {
			return Grade{
				Grade:  score.AbsenceIncompleteScore,
				Reason: GradeReasonIncorrectMultiplicity("", ref.Label.Text),
			}
		} else if correctnessPercentage < 0.99 {
			return Grade{
				Grade: correctnessPercentage * score.PresenceScore,
				Reason: GradeReasonIncorrectMultiplicity(
					fmt.Sprintf("partially (%.2f) ", correctnessPercentage),
					ref.Label.Text,
				),
			}
		} else {
			return Grade{Grade: score.PresenceScore, Reason: GRADE_REASON_EQUAL}
		}

	} else {
		c.LogErr("Grading of anything other than a multiplicity in the start/end of an edge is"+
			"not supported yet! (reference edge '%s', submission edge '%s')", refEdge.Id, subEdge.Id)
		return Grade{
			Grade:  0,
			Reason: GradeReasonUnsupported("non-multiplicity in edge end"),
		}
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

func syntacticSemanticMatch(c *context.Ctx, ref string, act string, semanticCertaintyThreshold float64, score RubricScoring) Grade {
	matchCmd := command.MatchStringCmd{
		Ref: ref,
		Act: act,
	}

	syntacticSimilarity := c.Queries.SyntacticMatch.Get("Syntactic similarity", matchCmd)
	semanticSimilarity := c.Queries.SemanticMatchSentenceTransformer.Get("Semantic similarity", matchCmd)

	if semanticSimilarity.Err != nil {
		c.LogErr("Could not calculate semantic similarity between '%s' and '%s': %s", ref, act, semanticSimilarity.Err.Error())
		semanticSimilarity.Score = 0
	}

	return formatPresentOrAbsenceGradeWithCorrectMsg(
		syntacticSimilarity < SYNTACTIC_DISTANCE_THRESHOLD ||
			semanticSimilarity.Score > semanticCertaintyThreshold,
		score,
		GRADE_REASON_EQUAL,
		GRADE_REASON_ABSENT_OR_INCORRECT,
	)
}

func formatPresentOrAbsenceGrade(correctGradeCondition bool, score RubricScoring, incorrectMsg GradeReason) Grade {
	return formatPresentOrAbsenceGradeWithCorrectMsg(correctGradeCondition, score, GRADE_REASON_PRESENT, incorrectMsg)
}

func formatPresentOrAbsenceGradeWithCorrectMsg(correctGradeCondition bool, score RubricScoring, correctMsg GradeReason, incorrectMsg GradeReason) Grade {
	res := Grade{
		Grade:  score.PresenceScore,
		Reason: correctMsg,
	}
	if !correctGradeCondition {
		res.Grade = score.AbsenceIncompleteScore
		res.Reason = incorrectMsg
	}
	return res
}
