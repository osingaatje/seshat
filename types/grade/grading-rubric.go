package result

// support for additive and deductive grading? 100% - deduction points or 0% + points?

type GradingRubric struct {
	// optional map of ILOs
	ILOs map[uint32]ILO

	Scoring DiagramMatchScores
}

type ILO struct {
	Id    uint32
	Name  string
	Descr string
}
type DiagramMatchScores struct {
	VertexScore          RubricScoring // presence/absence of vertex / node
	VertexAttributeScore RubricScoring // type of vertex

	AttributeScore     RubricScoring // field/method presence/absence
	AttributeTypeScore RubricScoring // wrong/right/extra attr. type

	AssociationScore             RubricScoring `json:"association"`              // presence/absence of association |v1|-----|v2|
	AssociationTypeScore         RubricScoring `json:"association_type"`         // association type |v1|<>-------|v2|
	AssociationMultiplicityScore RubricScoring `json:"association_multiplicity"` // association mult. |v1|*<>----1|v2|
	AssociationLabelScore        RubricScoring `json:"association_label"`        // association label |v1|*<>--text---1|v2|
}

type RubricScoring struct {
	PointsForPresence    float32 `json:"present"`     // if we have an element, how many points to award for presence
	PointsForAbsence     float32 `json:"absent"`      // if this element is missing, how many points to deduct
	PointsForSuperfluous float32 `json:"superfluous"` // if we have an extra element, how many points to deduct

	ILOWeight *ILOWeight // optional
}
type ILOWeight struct {
	ILO    uint32
	Weight float32 // 0.1, or 10 points, whatever the grader wants
}
