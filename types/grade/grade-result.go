package grade

import (
	"github.com/osingaatje/seshat/types/graph/intern"
)

type GradeResult struct {
	Reference  *intern.InternalGraph `json:"reference"`
	Submission *intern.InternalGraph `json:"submission"`
	Rubric     *GradeRubric          `json:"rubric"`
}
