package data

type ParseResult struct {
	Vertices map[uint64]ParsedVertex
	Edges    map[uint64]ParsedEdge
}

type ParsedVertex struct {
	Id         uint64
	Title      string                    // in UML, the classname
	Properties map[VertexProperty]string // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue    // raw values (in UML, the fields)

	// additional stuff for possible visualisation later on:
	Location Location2D
	Size     Location2D
}

type ParsedEdge struct {
	FromId         uint64
	ToId           uint64
	FromProperties map[EdgeEndProperty]string // things like multiplicity and arrow head style
	ToProperties   map[EdgeEndProperty]string

	Properties map[EdgeProperty]string // general properties such as edge label text etc.
}

// contains value along with optional properties
type ParsedValue struct {
	Value      string                   // raw value (i.e. the fieldValue in "fieldName: fieldValue"
	Properties map[ValueProperty]string // things like visibility etc.
}

// location on some grid or whatever
type Location2D struct {
	X float64
	Y float64
}

// TODO: add internal graph representation that removes some styling etc.
