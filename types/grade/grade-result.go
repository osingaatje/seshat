package grade

import (
	"fmt"

	"github.com/osingaatje/seshat/types/graph/shared"
)

type GradeReason string

const (
	GRADE_REASON_EQUAL               GradeReason = "equal to sample solution"
	GRADE_REASON_PRESENT             GradeReason = "present"
	GRADE_REASON_CORRECT             GradeReason = "correct"
	GRADE_REASON_ABSENT_OR_INCORRECT GradeReason = "is missing in the submission or is incorrect"
	GRADE_REASON_ABSENT              GradeReason = "is missing from the submission"
	GRADE_REASON_INCORRECT           GradeReason = "is incorrect"
	GRADE_REASON_SUPERFLUOUS         GradeReason = "superfluous: does not appear in reference solution"
)

func GradeReasonIncorrect(shouldHaveBeenName string) GradeReason {
	return GradeReason(fmt.Sprintf("is incorrect - should have been '%s'", shouldHaveBeenName))
}
func GradeReasonUnsupported(unsupportedThing string) GradeReason {
	return GradeReason(fmt.Sprintf("unsupported feature: %s", unsupportedThing))
}
func GradeReasonIncorrectMultiplicity(incorrectModifier string, shouldBe string) GradeReason {
	return GradeReason(fmt.Sprintf("%sincorrect multiplicity: should have been '%s'", incorrectModifier, shouldBe))
}

type GradeResult struct {
	FinalGrade  float64          `json:"final_grade"`
	Calculation GradeCalculation `json:"reason"`
}

type GradeCalculation struct {
	MissingReferenceVertices map[shared.VertexIdentifier]GradeResultVertex `json:"missing_reference_vertices"`
	VertexGrades             map[shared.VertexIdentifier]GradeResultVertex `json:"vertex_grades"`

	MissingReferenceEdges map[shared.EdgeIdentifier]GradeResultEdge `json:"missing_reference_edges"`
	EdgeGrades            map[shared.EdgeIdentifier]GradeResultEdge `json:"edge_grades"`
}

func NewGradeCalculation() GradeCalculation {
	return GradeCalculation{
		MissingReferenceVertices: map[shared.VertexIdentifier]GradeResultVertex{},
		VertexGrades:             map[shared.VertexIdentifier]GradeResultVertex{},

		MissingReferenceEdges: map[shared.EdgeIdentifier]GradeResultEdge{},
		EdgeGrades:            map[shared.EdgeIdentifier]GradeResultEdge{},
	}
}

type GradeResultVertex struct {
	PresenceScore   Grade                            `json:"presence"`
	TypeScore       *Grade                           `json:"type"`                 // optional
	TitleScore      *Grade                           `json:"title"`                //optional
	AttributeScores map[string]GradeResultVertexAttr `json:"attributes,omitempty"` //optional
}

func GradeResultVertexMissing(missingScore float64) GradeResultVertex {
	return GradeResultVertex{
		PresenceScore: Grade{
			Grade:  missingScore,
			Reason: GRADE_REASON_ABSENT,
		},
	}
}

func GradeResultVertexExtra(superfluousScore float64) GradeResultVertex {
	return GradeResultVertex{
		PresenceScore: Grade{
			Grade:  superfluousScore,
			Reason: GRADE_REASON_SUPERFLUOUS,
		},
	}
}

type GradeResultEdge struct {
	PresenceScore  Grade  `json:"presence"`
	EdgeTypeScore  *Grade `json:"line_style"`   // optional, whether the edge has the correct arrow styles and line style
	EdgeLabelScore *Grade `json:"label"`        // optional
	StartScore     *Grade `json:"vertex_start"` // optional
	EndScore       *Grade `json:"vertex_end"`   // optional
}

func GradeResultEdgeMissing(missingScore float64) GradeResultEdge {
	return GradeResultEdge{
		PresenceScore: Grade{
			Grade:  missingScore,
			Reason: GRADE_REASON_ABSENT,
		},
	}
}
func GradeResultEdgeExtra(superfluousScore float64) GradeResultEdge {
	return GradeResultEdge{
		PresenceScore: Grade{
			Grade:  superfluousScore,
			Reason: GRADE_REASON_SUPERFLUOUS,
		},
	}
}

type GradeResultVertexAttr struct {
	NameScore       Grade  `json:"name"`
	TypeScore       *Grade `json:"type,omitempty"`
	VisibilityScore *Grade `json:"visibility,omitempty"`
}

func GradeResultAttrMissing(missingScore float64) GradeResultVertexAttr {
	return GradeResultVertexAttr{
		NameScore: Grade{
			Grade:  missingScore,
			Reason: GRADE_REASON_ABSENT,
		},
	}
}
func GradeResultAttrExtra(superfluousScore float64) GradeResultVertexAttr {
	return GradeResultVertexAttr{
		NameScore: Grade{
			Grade:  superfluousScore,
			Reason: GRADE_REASON_SUPERFLUOUS,
		},
	}
}

type Grade struct {
	Grade  float64     `json:"grade"`
	Reason GradeReason `json:"reason"`
}
