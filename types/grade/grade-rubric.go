package result

// support for additive and deductive grading? 100% - deduction points or 0% + points?

type GradeRubric struct {
	// optional map of ILOs
	ILOs map[uint32]ILO `json:"ilos"`

	Scores DiagramMatchScores `json:"scores"`
}

type ILO struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Descr string `json:"description"`
}
type DiagramMatchScores struct {
	VertexScore                  RubricScoring `json:"vertex"`                   // presence/absence of vertex / node
	VertexTypeScore              RubricScoring `json:"vertex_type"`              // type of vertex
	VertexAttributeScore         RubricScoring `json:"vertex_attribute"`         // vertex attribute(s)
	VertexAttributeTypeScore     RubricScoring `json:"vertex_attribute_type"`    // vertex attribute type(s)
	AssociationScore             RubricScoring `json:"association"`              // presence/absence |v1|-----|v2|
	AssociationTypeScore         RubricScoring `json:"association_type"`         // association type |v1|<>-------|v2|
	AssociationMultiplicityScore RubricScoring `json:"association_multiplicity"` // association mult. |v1|*----1|v2|
	AssociationLabelScore        RubricScoring `json:"association_label"`        // association label |v1|--text---|v2|
}

type RubricScoring struct {
	PointsForPresence    float32 `json:"present"`     // if we have an element, how many points to award for presence
	PointsForAbsence     float32 `json:"absent"`      // if this element is missing, how many points to deduct
	PointsForSuperfluous float32 `json:"superfluous"` // if we have an extra element, how many points to deduct

	ILOWeight []ILOWeight `json:"ilo_weights,omitempty"` // optional
}
type ILOWeight struct {
	ILO    uint32  `json:"ilo_id"`
	Weight float32 `json:"weight"` // 0.1, or 10 points, whatever the grader wants
}
