package grade

import (
	"github.com/osingaatje/seshat/types/graph/shared"
)

type GradeResult struct {
	MissingVertices []shared.VertexIdentifier                     `json:"missing_vertices"`
	VertexGrades    map[shared.VertexIdentifier]GradeResultVertex `json:"vertex_grades"`

	MissingEdges []shared.EdgeIdentifier
	EdgeGrades   map[shared.EdgeIdentifier]GradeResultEdge `json:"edge_grades"`
}

func NewGradeResult() GradeResult {
	return GradeResult{
		MissingVertices: []shared.VertexIdentifier{},
		VertexGrades:    map[shared.VertexIdentifier]GradeResultVertex{},
		MissingEdges:    []shared.EdgeIdentifier{},
		EdgeGrades:      map[shared.EdgeIdentifier]GradeResultEdge{},
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
	MultiplicityScore Grade
}

type Grade struct {
	Grade  float64 `json:"grade"`
	Reason string  `json:"reason"`
}
