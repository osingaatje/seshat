package internalrep

import (
	"maps"

	. "github.com/osingaatje/seshat/types/graph/shared"
)

type InternalGraph struct {
	Vertices map[VertexIdentifier]*InternalVertex `json:"vertices"`
	Edges    map[EdgeIdentifier]*InternalEdge     `json:"edges"`
}

func NewGraph() InternalGraph {
	return InternalGraph{
		Vertices: map[VertexIdentifier]*InternalVertex{},
		Edges:    map[EdgeIdentifier]*InternalEdge{},
	}
}

func (g *InternalGraph) Copy() *InternalGraph {
	if g == nil {
		return nil
	}
	res := NewGraph()
	maps.Copy(res.Vertices, g.Vertices)
	maps.Copy(res.Edges, g.Edges)
	return &res
}

func (g *InternalGraph) VerticesConnect(v1 VertexIdentifier, v2 VertexIdentifier) ([]EdgeIdentifier, bool) {
	res := []EdgeIdentifier{}
	for eId, e := range g.Edges {
		if e.FromId != nil && e.ToId != nil && (*e.FromId == v1 || *e.FromId == v2) && (*e.ToId == v1 || *e.ToId == v2) {
			res = append(res, eId)
		}
	}
	return res, len(res) > 0
}

func (g *InternalGraph) ConnectedEdges(v VertexIdentifier) []EdgeIdentifier {
	res := []EdgeIdentifier{}
	for eId, e := range g.Edges {
		if (e.FromId != nil && *e.FromId == v) || (e.ToId != nil && *e.ToId == v) {
			res = append(res, eId)
		}
	}
	return res
}

func (g *InternalGraph) MergeConnectedEdges(edges []EdgeIdentifier, referenceGraph *InternalGraph) {
	if g == nil || referenceGraph == nil {
		panic("ONE OF THE GRAPHS WAS NIL WHEN MERGING CONNECTING EDGES!")
	}

	for _, eId := range edges {
		edge := referenceGraph.Edges[eId]
		g.Edges[eId] = edge
		if edge.FromId != nil {
			g.Vertices[*edge.FromId] = referenceGraph.Vertices[*edge.FromId]
		}
		if edge.ToId != nil {
			g.Vertices[*edge.ToId] = referenceGraph.Vertices[*edge.ToId]
		}
	}
}

func (g *InternalGraph) DeleteEdgesAndTheirVertices(edges []EdgeIdentifier) {
	if g == nil {
		panic("GRAPH NIL WHEN DELETING EDGES!")
	}

	for _, eId := range edges {
		edge, ok := g.Edges[eId]
		if !ok || edge == nil {
			continue
		}

		if edge.FromId != nil {
			delete(g.Vertices, *edge.FromId)
		}
		if edge.ToId != nil {
			delete(g.Vertices, *edge.ToId)
		}
		delete(g.Edges, eId)
	}
}

func (g *InternalGraph) Empty() bool {
	return len(g.Edges) == 0 && len(g.Vertices) == 0
}

type InternalVertex struct {
	Id         VertexIdentifier       `json:"id"`
	Title      string                 `json:"title"`      // in UML, the classname
	Properties VertexProperties       `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`

	// visual elements are omitted
}

type InternalEdge struct {
	FromId         *VertexIdentifier         `json:"fromId"`
	ToId           *VertexIdentifier         `json:"toId"`
	FromEdgeId     *EdgeIdentifier           `json:"fromEdgeId"`
	ToEdgeId       *EdgeIdentifier           `json:"toEdgeId"`
	FromProperties InternalEdgeEndProperties `json:"fromProperties"`  // things like multiplicity and arrow head style
	Label          *ParsedLabel              `json:"label,omitempty"` // for ex.: "teaches >"
	ToProperties   InternalEdgeEndProperties `json:"toProperties"`

	StyleProperties EdgeStyleProperties `json:"styleProperties"` // general properties such as edge label text etc.

	// visual properties omitted
}
