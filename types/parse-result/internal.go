package parseresult

import (
	"fmt"

	"github.com/osingaatje/seshat/helper"
)

type InternalGraphJSON struct {
	Vertices map[VertexIdentifier]*InternalVertex `json:"vertices"`
	Edges    map[string]*InternalEdge             `json:"edges"`
}

// Because we cannot marshal structs in map keys, we convert them to strings for json marshalling.
func (pr InternalGraph) MarshalJSON() ([]byte, error) {
	r := InternalGraphJSON{
		Vertices: pr.Vertices,
		Edges:    map[string]*InternalEdge{},
	}

	for _, edge := range pr.Edges {
		r.Edges[fmt.Sprintf("%d-%d", edge.FromId, edge.ToId)] = edge
	}
	return helper.MarshalJSON(r)
}

type InternalGraph struct {
	Vertices map[VertexIdentifier]*InternalVertex `json:"vertices"`
	Edges    map[EdgeIdentifier]*InternalEdge     `json:"edges"`
}

type InternalVertex struct {
	Id         VertexIdentifier       `json:"id"`
	Title      string                 `json:"title"`      // in UML, the classname
	Properties VertexProperties       `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`

	// visual elements are omitted
}

type InternalEdge struct {
	FromId         VertexIdentifier          `json:"fromId"`
	ToId           VertexIdentifier          `json:"toId"`
	FromProperties InternalEdgeEndProperties `json:"fromProperties"`  // things like multiplicity and arrow head style
	Label          *ParsedLabel              `json:"label,omitempty"` // for ex.: "teaches >"
	ToProperties   InternalEdgeEndProperties `json:"toProperties"`

	StyleProperties EdgeStyleProperties `json:"styleProperties"` // general properties such as edge label text etc.

	// visual properties omitted
}
