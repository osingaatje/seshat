package command

import (
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
)

type GradeCmd struct {
	Rubric            *GradeRubric
	ReferenceSolution *ParseResult
	Submission        *ParseResult
}
