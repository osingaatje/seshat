package command

import (
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/parse-result"
)

type GradeCmd struct {
	Rubric            *GradeRubric
	ReferenceSolution *InternalGraph
	Submission        *InternalGraph
}
