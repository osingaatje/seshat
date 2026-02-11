package data

// VERTICES
type VertexProperty string

var VertexPropertyAll []VertexProperty = []VertexProperty{
	ClassType,
	ClassVisibility,
}

const (
	ClassType       VertexProperty = "Class Type"
	ClassVisibility VertexProperty = "Visibility" // "public", "private", "protected", ...
	// TODO: add more if needed.

)

// EDGES
type EdgeEndProperty string

var EdgeEndPropertyAll []EdgeEndProperty = []EdgeEndProperty{
	ArrowStyle,
	EdgeMultiplicity,
}

const (
	ArrowStyle       EdgeEndProperty = "Arrow Style"
	EdgeMultiplicity EdgeEndProperty = "Multiplicity"
	// TODO: add more if needed.
)

type EdgeProperty string

var EdgePropertyAll []EdgeProperty = []EdgeProperty{
	LabelText,
}

const (
	LabelText EdgeProperty = "Label"
	// TODO: add more if needed.
)

// Properties of values
type ValueProperty string

var ValuePropertyAll []ValueProperty = []ValueProperty{
	ValueVisibility,
	ValueType,
}

const (
	ValueVisibility ValueProperty = "Visibility"
	ValueType       ValueProperty = "Type"
	// TODO: add more if needed.
)
