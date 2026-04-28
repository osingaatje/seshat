package intern

import (
	"fmt"

	"github.com/osingaatje/seshat/types/graph/shared"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

/*
 * We define the parse result visually as this:
 *     --------> x
 *     _____________________________________________
 * |  |
 * |  |
 *\/y |   .______             ._____
 *    |   | v1   | ---------- | v2 |
 *    |   --------            -----
 *
 *  . = position (X,Y)
 * 	width = X-offset (to the right) for a vertex
 *  height = Y-offset (to the bottom) for a vertex
 *
 */

type InternalGraph struct {
	Metadata shared.GraphMetadata `json:"metadata"`

	Vertices map[VertexIdentifier]*InternalVertex `json:"vertices"`
	Edges    map[EdgeIdentifier]*InternalEdge     `json:"edges"`
}

func (p *InternalGraph) Copy() *InternalGraph {
	if p == nil {
		return nil
	}

	res := NewParseResult()
	res.Metadata = p.Metadata.Copy()

	for k, v := range p.Vertices {
		res.Vertices[k] = v.Copy()
	}
	for k, e := range p.Edges {
		res.Edges[k] = e.Copy()
	}
	return res
}

func NewParseResult() *InternalGraph {
	res := InternalGraph{}
	res.Vertices = map[VertexIdentifier]*InternalVertex{}
	res.Edges = map[EdgeIdentifier]*InternalEdge{}
	return &res
}

func (g *InternalGraph) Empty() bool {
	return len(g.Edges) == 0 && len(g.Vertices) == 0
}

type InternalVertex struct {
	Id         VertexIdentifier       `json:"id"`         // unique ID to refer to from edges
	Title      string                 `json:"title"`      // in UML, the classname
	Properties VertexProperties       `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`     // raw values (in UML, the fields)

	VisualProperties VertexVisualProperties `json:"visual_properties"`
}

func (p *InternalVertex) Copy() *InternalVertex {
	if p == nil {
		return nil
	}

	res := &InternalVertex{
		Id:               p.Id,
		Title:            p.Title,
		Properties:       p.Properties,
		Values:           map[string]ParsedValue{}, // fill in manually to avoid having to DeepCopy
		VisualProperties: p.VisualProperties,
	}
	for k, v := range p.Values {
		res.Values[k] = v.Copy()
	}
	return res
}

func (p *InternalVertex) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("V(id=%d,title=%s)", p.Id, p.Title)
}

type InternalEdge struct {
	Id         EdgeIdentifier    `json:"id"`
	FromId     *VertexIdentifier `json:"fromId"`     // an edge may be 'floating', i.e. not connected
	ToId       *VertexIdentifier `json:"toId"`       // an edge may be 'floating', i.e. not connected
	FromEdgeId *EdgeIdentifier   `json:"fromEdgeId"` // an edge may be connected to another edge
	ToEdgeId   *EdgeIdentifier   `json:"toEdgeId"`   // an edge may be connected to another edge

	FromProperties EdgeEndProperties `json:"fromProperties"`  // things like multiplicity and arrow head style
	Label          *Label            `json:"label,omitempty"` // for ex.: "teaches >"
	ToProperties   EdgeEndProperties `json:"toProperties"`

	StyleProperties  EdgeStyleProperties  `json:"styleProperties"` // general properties such as edge label text etc.
	VisualProperties EdgeVisualProperties `json:"visualProperties"`
}

func (e *InternalEdge) Copy() *InternalEdge {
	if e == nil {
		return nil
	}

	res := &InternalEdge{
		Id:               e.Id,
		FromId:           e.FromId,
		ToId:             e.ToId,
		FromEdgeId:       e.FromEdgeId,
		ToEdgeId:         e.ToEdgeId,
		FromProperties:   e.FromProperties.Copy(),
		Label:            e.Label.Copy(),
		ToProperties:     e.ToProperties.Copy(),
		StyleProperties:  e.StyleProperties,
		VisualProperties: e.VisualProperties,
	}
	return res
}

func (e *InternalEdge) String() string {
	if e == nil {
		return "nil"
	}
	res := "E("
	if e.FromId != nil {
		res += fmt.Sprintf("fromId=%d,", *e.FromId)
	}
	if e.FromEdgeId != nil {
		res += fmt.Sprintf("fromEdgeId=%d,", *e.FromEdgeId)
	}
	if e.ToId != nil {
		res += fmt.Sprintf("toId=%d,", *e.ToId)
	}
	if e.ToEdgeId != nil {
		res += fmt.Sprintf("toEdgeId=%d,", *e.ToEdgeId)
	}
	res = res[:len(res)-1] // remove trailing comma
	res += ")"
	return res
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

// get all connected edges in a subgraph gsub
func (g *InternalGraph) ConnectedEdges(gsub *InternalGraph) []EdgeIdentifier {
	res := []EdgeIdentifier{}
	for eId, e := range g.Edges {
		if _, ok := gsub.Edges[eId]; ok { // skip already contained edges
			continue
		}

		if e.FromId != nil {
			if _, ok := gsub.Vertices[*e.FromId]; ok {
				res = append(res, eId)
				continue
			}
		}
		if e.ToId != nil {
			if _, ok := gsub.Vertices[*e.ToId]; ok {
				res = append(res, eId)
				continue
			}
		}

		if e.FromEdgeId != nil {
			if _, ok := gsub.Edges[*e.FromEdgeId]; ok {
				res = append(res, eId)
				continue
			}
		}
		if e.ToEdgeId != nil {
			if _, ok := gsub.Edges[*e.ToEdgeId]; ok {
				res = append(res, eId)
				continue
			}
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
