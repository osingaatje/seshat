package test

import (
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
	"github.com/osingaatje/seshat/types/graph/intern"
)

func TestGradeEdgeToEdgeGraphs(t *testing.T) {
	FILE1 := "./examples/similar/edge-to-edge/double-association.utml"
	FILE2 := "./examples/similar/edge-to-edge/double-association-wrong.utml"

	rubric := grade.NewGradeRubric()
	g1, g2, r, err := gradeStuff(FILE1, FILE2, rubric)
	if err != nil {
		t.Errorf("Error while grading: %s", err.Error())
		return
	}
	t.Error("TODO COMPLETE THIS TEST")
	assert.NotEqual(t, g1, g2, "bruh")
	assert.NotNil(t, r, "bruh2")
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
		c.LogErr("")
		return nil, nil, nil, fmt.Errorf("%sCould not interpret utml file(s) as parse result", logPrefix)
	}
	fixed1 = c.Queries.RepairDiagram.Get("Repair file 1", command.NewRepairCmdDefOpt(fixed1))
	fixed2 = c.Queries.RepairDiagram.Get("Repair file 2", command.NewRepairCmdDefOpt(fixed2))
	if fixed1 == nil || fixed2 == nil {
		return nil, nil, nil, fmt.Errorf("%sCould not repair one or both file(s)", logPrefix)
	}

	grade := c.Queries.GradeDiagram.Get("Grade file 1", command.GradeCmd{
		Rubric:            &rubric,
		ReferenceSolution: fixed1,
		Submission:        fixed2,
	})
	if grade == nil {
		return fixed1, fixed2, nil, fmt.Errorf("%sFailed to grade submission", logPrefix)
	}

	return fixed1, fixed2, grade, nil
}
