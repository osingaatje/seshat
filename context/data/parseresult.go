package data

type ParseResult struct {
	Vertices map[uint64]*ParsedVertex `json:"vertices"`
	Edges    map[uint64]*ParsedEdge   `json:"edges"`
}

func NewParseResult() *ParseResult {
	res := ParseResult{}
	res.Vertices = map[uint64]*ParsedVertex{}
	res.Edges = map[uint64]*ParsedEdge{}
	return &res
}

type ParsedVertex struct {
	Id         uint64                    `json:"id"`
	Title      string                    `json:"title"`      // in UML, the classname
	Properties map[VertexProperty]string `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue    `json:"values"`     // raw values (in UML, the fields)

	// additional stuff for possible visualisation later on:
	Location Location2D `json:"location"`
	Size     Location2D `json:"size"`
}

type ParsedEdge struct {
	FromId         uint64                     `json:"fromId"`
	ToId           uint64                     `json:"toId"`
	FromProperties map[EdgeEndProperty]string `json:"fromProperties"` // things like multiplicity and arrow head style
	Label          ParsedLabel                `json:"label"`          // for ex.: "teaches >"
	ToProperties   map[EdgeEndProperty]string `json:"toProperties"`

	Properties map[EdgeProperty]string `json:"properties"` // general properties such as edge label text etc.
}

// contains value along with optional properties
type ParsedValue struct {
	Value      string                   `json:"value"`      // raw value (i.e. the fieldValue in "fieldName: fieldValue"
	Properties map[ValueProperty]string `json:"properties"` // things like visibility etc.
}

type ParsedLabel struct {
	Text     string     `json:"text"`
	Location Location2D `json:"location"`
}

// location on some grid or whatever
type Location2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
