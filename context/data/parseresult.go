package data

import (
	. "github.com/osingaatje/seshat/context/data/parse-result-datatypes"
)

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
	Id         uint64                 `json:"id"`         // unique ID to refer to from edges
	Title      string                 `json:"title"`      // in UML, the classname
	Properties map[VertexProperty]any `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`     // raw values (in UML, the fields)

	// additional stuff for possible visualisation later on:
	Location Vector2D `json:"location"`
	Size     Vector2D `json:"size"`
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
	Value      string                `json:"value"`      // raw value (i.e. the fieldValue in "fieldName: fieldValue"
	Properties map[ValueProperty]any `json:"properties"` // things like visibility etc.
}

// like in the UML text along an edge. Also contains a location so we can route the edge along this label
type ParsedLabel struct {
	Text     string   `json:"text"`
	Location Vector2D `json:"location"`
}

// location on some grid or whatever
type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
