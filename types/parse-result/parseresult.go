package parseresult

import (
	"fmt"
	"strings"

	"github.com/osingaatje/seshat/helper"
	. "github.com/osingaatje/seshat/types/generic"
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

// Because we cannot marshal structs in map keys, we convert them to strings for json marshalling.
func (pr ParseResult) MarshalJSON() ([]byte, error) {
	r := ParseResultJSON{
		Vertices: pr.Vertices,
		Edges:    map[string]*ParsedEdge{},
	}

	for _, edge := range pr.Edges {
		r.Edges[fmt.Sprintf("%d-%d", edge.FromId, edge.ToId)] = edge
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

type VertexIdentifier uint32
type EdgeIdentifier uint64

func NewEdgeIdentifier(vId1 VertexIdentifier, vId2 VertexIdentifier) EdgeIdentifier {
	return EdgeIdentifier(uint64(vId1)<<32 | uint64(vId2))
}
func (e EdgeIdentifier) New(vId1 VertexIdentifier, vId2 VertexIdentifier) EdgeIdentifier {
	return NewEdgeIdentifier(vId1, vId2)
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
		res.Vertices[id] = v.ToInternal()
	}
	for id, e := range p.Edges {
		iE, err := e.ToInternal()
		if err != nil {
			errors = append(errors, err.Error())
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
func (p *ParsedVertex) ToInternal() *InternalVertex {
	if p == nil {
		return nil
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
	return res
}

type ParsedEdge struct {
	FromId         VertexIdentifier  `json:"fromId"`
	ToId           VertexIdentifier  `json:"toId"`
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
		FromId:          e.FromId,
		ToId:            e.ToId,
		FromProperties:  e.FromProperties.Copy(),
		Label:           e.Label.Copy(),
		ToProperties:    e.ToProperties.Copy(),
		StyleProperties: e.StyleProperties,
	}
	return res
}

func (e *ParsedEdge) ToInternal() (*InternalEdge, error) {
	res := &InternalEdge{
		FromId:          e.FromId,
		ToId:            e.ToId,
		Label:           e.Label.Copy(),
		StyleProperties: e.StyleProperties,
		// visual props omitted

		// FromProps, ToProps added below
	}

	fromP, okF := e.FromProperties.ToInternal()
	toP, okT := e.ToProperties.ToInternal()
	if !okF {
		return nil, fmt.Errorf("Could not parse edge start multplicity (%d-%d)", e.FromId, e.ToId)
	}
	if !okT {
		return nil, fmt.Errorf("Could not parse edge end multiplicity (%d-%d)", e.FromId, e.ToId)
	}
	res.FromProperties = fromP
	res.ToProperties = toP
	return res, nil
}

// contains value along with optional properties
type ParsedValue struct {
	Value      string          `json:"value"`      // raw value (i.e. the fieldValue in "fieldName: fieldValue"
	Properties ValueProperties `json:"properties"` // things like visibility etc.
}

func (p ParsedValue) Copy() ParsedValue {
	return ParsedValue{
		Value:      p.Value,
		Properties: p.Properties.Copy(),
	}
}

// like in the UML text along an edge. Also contains a location so we can route the edge along this label
type ParsedLabel struct {
	Text     string   `json:"text"`
	Location Vector2D `json:"location"`
}

func (l *ParsedLabel) Copy() *ParsedLabel {
	if l == nil {
		return nil
	}
	return &ParsedLabel{
		Text:     l.Text,
		Location: l.Location,
	}
}
