package result

// support for additive and deductive grading? 100% - deduction points or 0% + points?

type GradingRubric struct {
	// optional map of ILOs
	ILOs map[uint32]ILO

	GeneralRubric GenericRubric
	Rubrics       []GradingRubricElement
}

type ILO struct {
	Id    uint32
	Name  string
	Descr string
}
type GenericRubric struct {
	VertexScore          RubricScoring // presence/absence of vertex / node
	VertexAttributeScore RubricScoring // type of vertex

	AttributeScore     RubricScoring // field/method presence/absence
	AttributeTypeScore RubricScoring // wrong/right/extra attr. type

	AssociationLabelScore RubricScoring // presence/absence of correct label

	AssociationScore             RubricScoring // <-> presence/absence
	AssociationTypeScore         RubricScoring // correct/wrong type
	AssociationMultiplicityScore RubricScoring // presence/absence of
}

type RubricScoring struct {
	PointsForPresence    float32 // if we have an element, how many points to award / how many points to deduct for absence
	PointsForAbsence     float32
	PointsForSuperfluous float32 // if we have an extra element, how many points to deduct
}

// examples:
type GradingRubricElement struct {
	Rule GradingRule
	Args []any
}

type GradingRule string

const (
	VertexMustExist GradingRule = "Vertex '%s' must exist"
)

func (r GradingRule) Format(args []string)
