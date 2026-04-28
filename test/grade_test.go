package test

import (
	"testing"

	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
)

func TestGradeEdgeToEdgeGraphs(t *testing.T) {
	FILE1 := "./examples/similar/edge-to-edge/double-association.utml"
	FILE2 := "./examples/similar/edge-to-edge/double-association-wrong.utml"

	r := gradeStuff(FILE1, FILE2)
	if r == nil {
		t.Error("Grade result was nil. :(")
		return
	}
	t.Error("TODO COMPLETE THIS TEST")
}

func TestGradeSimilarDiagWrongMultiplicities(t *testing.T) {
	FILE1 := "./examples/similar/diff-assoc/simple-assocation-2-to-1-5.utml"
	FILE2 := "./examples/similar/diff-assoc/simple-assocation-one-to-many.utml"

	r := gradeStuff(FILE1, FILE2)
	if r == nil {
		t.Error("Grade result was nil. :(")
		return
	}
	t.Error("TODO COMPLETE THIS TEST")
}

func gradeStuff(filename1 string, filename2 string) *grade.GradeResult {
	c := driver.NewContext()

	utml1 := c.Queries.ParseUTML.Get("Parse UTML file 1", filename1)
	utml2 := c.Queries.ParseUTML.Get("Parse UTML file 1", filename2)
	if utml1 == nil || utml2 == nil {
		c.LogErr("Could not parse UTML file(s)")
		return nil
	}
	fixed1 := c.Queries.ParseUTMLToParseRes.Get("Interpret file1", utml1)
	fixed2 := c.Queries.ParseUTMLToParseRes.Get("Interpret file2", utml2)
	if fixed1 == nil || fixed2 == nil {
		c.LogErr("Could not interpret utml file(s) as parse result")
		return nil
	}
	fixed1 = c.Queries.RepairDiagram.Get("Repair file 1", command.NewRepairCmdDefOpt(fixed1))
	fixed2 = c.Queries.RepairDiagram.Get("Repair file 2", command.NewRepairCmdDefOpt(fixed2))
	if fixed1 == nil || fixed2 == nil {
		c.LogErr("Could not repair ")
		return nil
	}

	rubric := grade.NewGradeRubric()
	grade := c.Queries.GradeDiagram.Get("Grade file 1", command.GradeCmd{
		Rubric:            &rubric,
		ReferenceSolution: fixed1,
		Submission:        fixed2,
	})

	return grade
}
