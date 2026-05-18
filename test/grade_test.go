package test

import (
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
	"github.com/osingaatje/seshat/types/graph/intern"
	"github.com/osingaatje/seshat/types/graph/shared"
)

func TestGradeEdgeToEdgeGraphs(t *testing.T) {
	/*
	 * These diagrams have the same vertices, but one wrong edge (which is also connected to another edge)
	 */
	FILE1 := "./examples/similar/edge-to-edge/double-association.utml"
	FILE2 := "./examples/similar/edge-to-edge/double-association-wrong.utml"

	rubric := grade.NewGradeRubric()
	_, _, r, err := gradeStuff(FILE1, FILE2, rubric)
	if err != nil {
		t.Errorf("Error while grading: %s", err.Error())
		return
	}

	assert.Equal(t, 1, len(r.MissingReferenceEdges), "Should be one missing edge")
	assert.Equal(t, 3, len(r.EdgeGrades), "Should be three edge grades (2 correct, 1 superfluous)")
	assert.Equal(t, 1, len(
		helper.FilterMap(
			r.EdgeGrades,
			func(id shared.EdgeIdentifier, e grade.GradeResultEdge) bool {
				return e.PresenceScore.Grade == rubric.Scores.EdgeScore.SuperfluousScore
			},
		)), "Should be one 'superfluous' edge out of 3 total edges")

	for _, subVGrade := range r.VertexGrades {
		assert.Equal(t, rubric.Scores.VertexScore.PresenceScore, subVGrade.PresenceScore.Grade, "All vertices should be matched")
		assert.Equal(t, grade.GRADE_REASON_PRESENT, subVGrade.PresenceScore.Reason, "All vertices should receive a 'present' reason")
	}
	for _, subEGrade := range r.EdgeGrades {
		if subEGrade.PresenceScore.Grade == rubric.Scores.EdgeScore.PresenceScore {
			// present edge should have a match on all rubrics
			assert.Equal(t, rubric.Scores.EdgeScore.PresenceScore, subEGrade.PresenceScore.Grade, "Edge should be present")
			assert.Equal(t, grade.GRADE_REASON_PRESENT, subEGrade.PresenceScore.Reason, "Edge should get 'present' reason")

			for name, gres := range map[string]*grade.Grade{
				"edgeType":  subEGrade.EdgeTypeScore,
				"edgeLabel": subEGrade.EdgeLabelScore,
				"edgeStart": subEGrade.StartScore,
				"edgeEnd":   subEGrade.EndScore} {
				assert.Equal(t, rubric.Scores.EdgeScore.PresenceScore, gres.Grade, "Edge subgrade %s should be present", name)
				assert.Equal(t, grade.GRADE_REASON_EQUAL, gres.Reason, "Edge subgrade %s should get 'equal' reason", name)
			}

		} else {
			// should be the one superfluous edge
			assert.Equal(t, rubric.Scores.EdgeScore.SuperfluousScore, subEGrade.PresenceScore.Grade, "Superfluous should get 'superfluous' score")
			assert.Equal(t, grade.GRADE_REASON_SUPERFLUOUS, subEGrade.PresenceScore.Reason, "Edge should get 'superfluous' reason")
		}
	}
	for _, subEGrade := range r.MissingReferenceEdges {
		assert.Equal(t, rubric.Scores.EdgeScore.AbsenceIncompleteScore, subEGrade.PresenceScore.Grade, "Absent edge should get 'absent' score")
		assert.Equal(t, grade.GRADE_REASON_ABSENT, subEGrade.PresenceScore.Reason, "Absent edge should get 'absent' reason")
	}
}

func TestGradeSimilarDiagWrongMultiplicities(t *testing.T) {
	FILE1 := "./examples/similar/diff-assoc/simple-assocation-2-to-1-5.utml"
	FILE2 := "./examples/similar/diff-assoc/simple-assocation-one-to-many.utml"

	rubric := grade.NewGradeRubric()
	_, g2, r, err := gradeStuff(FILE1, FILE2, rubric)
	if err != nil {
		t.Errorf("Error while grading: %s", err.Error())
		return
	}

	assert.Equal(t, 1, len(g2.Edges), "Should be one edge in this example")

	// all vertices and edges should be present:
	assert.Equal(t, 0, len(r.MissingReferenceEdges), "Should not be any missing edges!")
	assert.Equal(t, 0, len(r.MissingReferenceVertices), "Should not be any missing vertices!")

	for subVId, _ := range g2.Vertices {
		subVGrade, ok := r.VertexGrades[subVId]
		assert.True(t, ok, "A submission vertex was not in the grading result!")
		assert.Equal(t, subVGrade.PresenceScore.Grade, rubric.Scores.VertexScore.PresenceScore, "Vertex did not get presence grade?")
		assert.Equal(t, subVGrade.PresenceScore.Reason, grade.GRADE_REASON_PRESENT, "Vertex did not get 'present' reason")
	}

	for subEId, _ := range g2.Edges {
		subEGrade, ok := r.EdgeGrades[subEId]
		assert.True(t, ok, "A submission edge was not in the grading result!")
		assert.Equal(t, subEGrade.PresenceScore.Grade, rubric.Scores.EdgeScore.PresenceScore, "Edge should have gotten presence score")
		assert.Equal(t, subEGrade.PresenceScore.Reason, grade.GRADE_REASON_PRESENT, "Edge did not get 'present' reason")

		assert.NotNil(t, subEGrade.StartScore)
		assert.Equal(t, subEGrade.StartScore.Grade, rubric.Scores.EdgeEndLabelScore.AbsenceIncompleteScore, "Muliplicity should be graded as 'wrong'")
		assert.NotNil(t, subEGrade.EndScore)
		assert.Equal(t, subEGrade.EndScore.Grade, rubric.Scores.EdgeEndLabelScore.AbsenceIncompleteScore, "Muliplicity should be graded as 'wrong'")
	}
}

func gradeStuff(filename1 string, filename2 string, rubric grade.GradeRubric) (
	g1 *intern.InternalGraph,
	g2 *intern.InternalGraph,
	gradeRes *grade.GradeCalculation,
	err error,
) {
	logPrefix := fmt.Sprintf("Grade '%s'<->'%s': ", filename1, filename2)
	c := driver.NewContext()

	utml1 := c.Queries.ParseUTML.Get("Parse UTML file 1", filename1)
	utml2 := c.Queries.ParseUTML.Get("Parse UTML file 1", filename2)
	if utml1 == nil || utml2 == nil {
		return nil, nil, nil, fmt.Errorf("%sCould not parse one or both UTML file(s)", logPrefix)
	}
	fixed1 := c.Queries.ParseUTMLToParseRes.Get("Interpret file1", utml1)
	fixed2 := c.Queries.ParseUTMLToParseRes.Get("Interpret file2", utml2)
	if fixed1 == nil || fixed2 == nil {
		return nil, nil, nil, fmt.Errorf("%sCould not interpret utml file(s) as parse result", logPrefix)
	}
	repair1 := c.Queries.RepairDiagram.Get("Repair file 1", command.NewRepairCmdDefOpt(fixed1))
	repair2 := c.Queries.RepairDiagram.Get("Repair file 2", command.NewRepairCmdDefOpt(fixed2))
	if len(repair1.Errors) > 0 || len(repair2.Errors) > 0 || repair1.Diagram == nil || repair2.Diagram == nil {
		return nil, nil, nil, fmt.Errorf("%sCould not repair one or both file(s): file 1 err: \"%s\" --- file 2 err: \"%s\"", logPrefix, repair1.Error(), repair2.Error())
	}

	grade := c.Queries.GradeDiagram.Get("Grade file 1", command.GradeCmd{
		Rubric:            &rubric,
		ReferenceSolution: repair1.Diagram,
		Submission:        repair2.Diagram,
	})
	if grade == nil {
		return fixed1, fixed2, nil, fmt.Errorf("%sFailed to grade submission", logPrefix)
	}

	return fixed1, fixed2, grade, nil
}
