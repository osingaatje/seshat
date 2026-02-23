package types

import (
	. "github.com/osingaatje/seshat/types/parse-result-datatypes"
)

type ParseResult struct {
	Vertices map[VertexIdentifier]*ParsedVertex `json:"vertices"`
	Edges    map[EdgeIdentifier]*ParsedEdge     `json:"edges"`
}

type VertexIdentifier uint64
type EdgeIdentifier struct {
	FromId VertexIdentifier
	ToId   VertexIdentifier
}

func NewParseResult() *ParseResult {
	res := ParseResult{}
	res.Vertices = map[VertexIdentifier]*ParsedVertex{}
	res.Edges = map[EdgeIdentifier]*ParsedEdge{}
	return &res
}

type ParsedVertex struct {
	Id         VertexIdentifier       `json:"id"`         // unique ID to refer to from edges
	Title      string                 `json:"title"`      // in UML, the classname
	Properties map[VertexProperty]any `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`     // raw values (in UML, the fields)

	// additional stuff for possible visualisation later on:
	Location Vector2D `json:"location"`
	Size     Vector2D `json:"size"`
}

type ParsedEdge struct {
	FromId         VertexIdentifier        `json:"fromId"`
	ToId           VertexIdentifier        `json:"toId"`
	FromProperties map[EdgeEndProperty]any `json:"fromProperties"` // things like multiplicity and arrow head style
	Label          ParsedLabel             `json:"label"`          // for ex.: "teaches >"
	ToProperties   map[EdgeEndProperty]any `json:"toProperties"`

	StyleProperties map[EdgeStyleProperty]any `json:"styleProperties"` // general properties such as edge label text etc.
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
