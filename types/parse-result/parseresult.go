package parseresult

import (
	"fmt"

	"github.com/osingaatje/seshat/helper"
)

// Because we cannot marshal structs in map keys, we convert them to strings for json marshalling.
func (pr ParseResult) MarshalJSON() ([]byte, error) {
	r := ParseResultJSON{
		Vertices: pr.Vertices,
		Edges:    map[string]*ParsedEdge{},
	}

	for id, edge := range pr.Edges {
		r.Edges[fmt.Sprintf("%d-%d", id.FromId, id.ToId)] = edge
	}
	return helper.MarshalJSON(r)
}

type ParseResultJSON struct {
	Vertices map[VertexIdentifier]*ParsedVertex `json:"vertices"`
	Edges    map[string]*ParsedEdge             `json:"edges"`
}

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
	Properties VertexProperties       `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`     // raw values (in UML, the fields)

	VisualProperties VertexVisualProperties `json:"visual_properties"`
}

type ParsedEdge struct {
	FromId         VertexIdentifier  `json:"fromId"`
	ToId           VertexIdentifier  `json:"toId"`
	FromProperties EdgeEndProperties `json:"fromProperties"` // things like multiplicity and arrow head style
	Label          ParsedLabel       `json:"label"`          // for ex.: "teaches >"
	ToProperties   EdgeEndProperties `json:"toProperties"`

	StyleProperties EdgeStyleProperties `json:"styleProperties"` // general properties such as edge label text etc.
}

// contains value along with optional properties
type ParsedValue struct {
	Value      string          `json:"value"`      // raw value (i.e. the fieldValue in "fieldName: fieldValue"
	Properties ValueProperties `json:"properties"` // things like visibility etc.
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
