package command

import (
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/intern"
)

type GradeCmd struct {
	Rubric            *GradeRubric
	ReferenceSolution *InternalGraph
	Submission        *InternalGraph
}
