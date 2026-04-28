package grade

import (
	"github.com/osingaatje/seshat/types/graph/shared"
)

type GradeReason string

const (
	GRADE_REASON_PRESENT     GradeReason = "present"
	GRADE_REASON_ABSENT      GradeReason = "appears in reference solution but is not in submission"
	GRADE_REASON_SUPERFLUOUS GradeReason = "does not appear in reference solution"
)

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
	PresenceScore   Grade `json:"presence"`
	TypeScore       Grade `json:"type"`
	TitleScore      Grade `json:"title"`
	AttributeScores map[string]AttributeScore
}

type GradeResultEdge struct {
	PresenceScore Grade `json:"presence"`
	TypeScore     Grade `json:"type"`

	EdgeStyleScore EdgeStyle    `json:"line_style"`
	EdgeLabelScore Grade        `json:"label"`
	StartScore     EdgeEndScore `json:"vertex_start"`
	EndScore       EdgeEndScore `json:"vertex_end"`
}

type AttributeScore struct {
	NameScore       Grade `json:"name"`
	TypeScore       Grade `json:"type"`
	VisibilityScore Grade `json:"visibility"`
}

type EdgeStyle struct {
	LineTypeScore Grade `json:"line_type"`
}

type EdgeEndScore struct {
	MultiplicityScore Grade `json:"multiplicity"`
}

type Grade struct {
	Grade  float64     `json:"grade"`
	Reason GradeReason `json:"reason"`
}
