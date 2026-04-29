package grade

// support for additive and deductive grading? 100% - deduction points or 0% + points?

type GradeRubric struct {
	// optional map of ILOs
	ILOs map[uint32]ILO `json:"ilos"`

	GraderConfig GraderConfig `json:"grader_config"`

	Scores DiagramMatchScores `json:"scores"`
}

func NewGradeRubric() GradeRubric {
	return GradeRubric{
		ILOs:         map[uint32]ILO{},
		GraderConfig: NewGraderConfig(),
		Scores:       NewDiagramMatchScores(),
	}
}

type ILO struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Descr string `json:"description"`
}

type GraderConfig struct {
	ClassContentSimilarity float64 `json:"semantic_certainty"` // how much do classes need to look alike in order to be matched?
}

func NewGraderConfig() GraderConfig {
	return GraderConfig{
		ClassContentSimilarity: 0.75, // TODO ACTUALLY USE THIS
	}
}

type DiagramMatchScores struct {
	VertexScore                    RubricScoring `json:"vertex"` // presence/absence of vertex / node
	VertexTitleScore               RubricScoring `json:"vertex_title"`
	VertexTypeScore                RubricScoring `json:"vertex_type"`
	VertexVisibilityScore          RubricScoring `json:"vertex_visibility"`
	VertexAttributeScore           RubricScoring `json:"vertex_attribute"`
	VertexAttributeTypeScore       RubricScoring `json:"vertex_attribute_type"`
	VertexAttributeVisibilityScore RubricScoring `json:"vertex_attribute_visibility"`
	EdgeScore                      RubricScoring `json:"edge"`                   // presence/absence |v1|-----|v2|
	EdgeTypeScore                  RubricScoring `json:"edge_type"`              // association type |v1|<>-------|v2|
	EdgeEndLabelScore              RubricScoring `json:"edge_association_label"` // association mult. |v1|*----1|v2|
	EdgeLabelScore                 RubricScoring `json:"edge_middle_label"`      // association label |v1|--text---|v2|
}

func NewDiagramMatchScores() DiagramMatchScores {
	return DiagramMatchScores{
		VertexScore:                    NewRubricScore(),
		VertexTitleScore:               NewRubricScore(),
		VertexTypeScore:                NewRubricScore(),
		VertexVisibilityScore:          NewRubricScore(),
		VertexAttributeScore:           NewRubricScore(),
		VertexAttributeTypeScore:       NewRubricScore(),
		VertexAttributeVisibilityScore: NewRubricScore(),
		EdgeScore:                      NewRubricScore(),
		EdgeTypeScore:                  NewRubricScore(),
		EdgeEndLabelScore:              NewRubricScore(),
		EdgeLabelScore:                 NewRubricScore(),
	}
}

type RubricScoring struct {
	PresenceScore          float64 `json:"present"`          // if we have an element, how many points to award for presence
	AbsenceIncompleteScore float64 `json:"absent/incorrect"` // if this element is missing, how many points to deduct
	SuperfluousScore       float64 `json:"superfluous"`      // if we have an extra element, how many points to deduct

	ILOWeight []ILOWeight `json:"ilo_weights,omitempty"` // optional
}

func NewRubricScore() RubricScoring {
	return RubricScoring{
		PresenceScore:          +1,
		AbsenceIncompleteScore: -1,
		SuperfluousScore:       -.5,
		ILOWeight:              []ILOWeight{},
	}
}

type ILOWeight struct {
	ILO    uint32  `json:"ilo_id"`
	Weight float64 `json:"weight"` // 0.1, or 10 points, whatever the grader wants
}
