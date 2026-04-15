package parseresult

import (
	"fmt"
	"slices"
	"strings"

	. "github.com/osingaatje/seshat/types/graph/internal-rep"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

var SKIPPED_TYPES []string = []string{"CommentNode"}

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

type ParseResult struct {
	Vertices map[VertexIdentifier]*ParsedVertex `json:"vertices"`
	Edges    map[EdgeIdentifier]*ParsedEdge     `json:"edges"`
}

func (p *ParseResult) Copy() *ParseResult {
	if p == nil {
		return nil
	}

	res := NewParseResult()
	for k, v := range p.Vertices {
		res.Vertices[k] = v.Copy()
	}
	for k, e := range p.Edges {
		res.Edges[k] = e.Copy()
	}
	return res
}

func NewParseResult() *ParseResult {
	res := ParseResult{}
	res.Vertices = map[VertexIdentifier]*ParsedVertex{}
	res.Edges = map[EdgeIdentifier]*ParsedEdge{}
	return &res
}

func (p *ParseResult) ToInternal() (*InternalGraph, error) {
	if p == nil {
		return nil, fmt.Errorf("Parse Result was nil!")
	}

	res := &InternalGraph{
		Vertices: map[VertexIdentifier]*InternalVertex{},
		Edges:    map[EdgeIdentifier]*InternalEdge{},
	}
	errors := []string{}

	for id, v := range p.Vertices {
		iV, err := v.ToInternal()
		if err != nil {
			newErrMsg := fmt.Sprintf("Failed converting vertex %s to internal repr.: %s", v.String(), err.Error())
			errors = append(errors, newErrMsg)
			continue
		}
		res.Vertices[id] = iV
	}
	for id, e := range p.Edges {
		iE, err := e.ToInternal()
		if err != nil {
			newErrMsg := fmt.Sprintf("Failed converting edge %s to internal repr.: %s", e.String(), err.Error())
			errors = append(errors, newErrMsg)
		}
		res.Edges[id] = iE
	}

	var err error = nil
	if len(errors) > 0 {
		err = fmt.Errorf(
			"Errors happended while converting to internal result: [%s]",
			strings.Join(errors, ","),
		)
	}

	return res, err
}

type ParsedVertex struct {
	Id         VertexIdentifier       `json:"id"`         // unique ID to refer to from edges
	Title      string                 `json:"title"`      // in UML, the classname
	Properties VertexProperties       `json:"properties"` // things like the visibility, inheritance properties etc.
	Values     map[string]ParsedValue `json:"values"`     // raw values (in UML, the fields)

	VisualProperties VertexVisualProperties `json:"visual_properties"`
}

func (p *ParsedVertex) Copy() *ParsedVertex {
	if p == nil {
		return nil
	}

	res := &ParsedVertex{
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
func (p *ParsedVertex) ToInternal() (v *InternalVertex, err error) {
	if p == nil {
		return nil, fmt.Errorf("p is nil")
	}

	if slices.Contains(SKIPPED_TYPES, p.Properties.Type) {
		return nil, fmt.Errorf("Vertex %s will be skipped (contains skippable vertex type '%s')", p.String(), p.Properties.Type)
	}

	res := &InternalVertex{
		Id:         p.Id,
		Title:      p.Title,
		Properties: p.Properties,
		Values:     map[string]ParsedValue{},
	}
	for k, v := range p.Values {
		res.Values[k] = v.Copy()
	}
	return res, nil
}
func (p *ParsedVertex) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("V(id=%d,title=%s)", p.Id, p.Title)
}

type ParsedEdge struct {
	FromId     *VertexIdentifier `json:"fromId"`     // an edge may be 'floating', i.e. not connected
	ToId       *VertexIdentifier `json:"toId"`       // an edge may be 'floating', i.e. not connected
	FromEdgeId *EdgeIdentifier   `json:"fromEdgeId"` // an edge may be connected to another edge
	ToEdgeId   *EdgeIdentifier   `json:"toEdgeId"`   // an edge may be connected to another edge

	FromProperties EdgeEndProperties `json:"fromProperties"`  // things like multiplicity and arrow head style
	Label          *ParsedLabel      `json:"label,omitempty"` // for ex.: "teaches >"
	ToProperties   EdgeEndProperties `json:"toProperties"`

	StyleProperties  EdgeStyleProperties  `json:"styleProperties"` // general properties such as edge label text etc.
	VisualProperties EdgeVisualProperties `json:"visualProperties"`
}

func (e *ParsedEdge) Copy() *ParsedEdge {
	if e == nil {
		return nil
	}

	res := &ParsedEdge{
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

func (e *ParsedEdge) ToInternal() (*InternalEdge, error) {
	res := &InternalEdge{
		FromId:          e.FromId,
		ToId:            e.ToId,
		FromEdgeId:      e.FromEdgeId,
		ToEdgeId:        e.ToEdgeId,
		Label:           e.Label.Copy(),
		StyleProperties: e.StyleProperties,
		// visual props omitted

		// FromProps, ToProps added below
	}

	fromP, errF := e.FromProperties.ToInternal()
	toP, errT := e.ToProperties.ToInternal()
	if errF != nil {
		return nil, fmt.Errorf("Could not parse start properties: %s", errF.Error())
	}
	if errT != nil {
		return nil, fmt.Errorf("Could not parse end properties: %s", errT.Error())
	}
	res.FromProperties = fromP
	res.ToProperties = toP
	return res, nil
}
func (e *ParsedEdge) String() string {
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
