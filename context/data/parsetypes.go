package data

// VERTICES
type VertexProperty string

const (
	ClassType       VertexProperty = "Class type"
	ClassVisibility VertexProperty = "Visibility" // "public", "private", "protected", ...
	// TODO: add more if needed.
)

// EDGES
type EdgeEndProperty string

const (
	ArrowStyle       EdgeProperty    = "Arrow style"
	EdgeMultiplicity EdgeEndProperty = "Multiplicity"
	// TODO: add more if needed.
)

type EdgeProperty string

const (
	LabelText EdgeProperty = "Label"
	// TODO: add more if needed.
)

// Properties of values
type ValueProperty string

const (
	ValueVisibility ValueProperty = "Visibility"
	ValueType       ValueProperty = "Type"
	// TODO: add more if needed.
)
