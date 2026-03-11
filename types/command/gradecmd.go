package command

import (
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/internal-rep"
)

type GradeCmd struct {
	Rubric            *GradeRubric
	ReferenceSolution *InternalGraph
	Submission        *InternalGraph
}
