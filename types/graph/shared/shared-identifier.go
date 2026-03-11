package shared

import (
	. "github.com/osingaatje/seshat/types/generic"
)

type VertexIdentifier uint32
type EdgeIdentifier uint64

func NewEdgeIdentifier(vId1 VertexIdentifier, vId2 VertexIdentifier) EdgeIdentifier {
	return EdgeIdentifier(uint64(vId1)<<32 | uint64(vId2))
}
func (e EdgeIdentifier) New(vId1 VertexIdentifier, vId2 VertexIdentifier) EdgeIdentifier {
	return NewEdgeIdentifier(vId1, vId2)
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

// Properties of values, e.g.: "field" : Visibility=private, Type=string
type ValueProperties struct {
	Visibility ValuePropVisibilityVar `json:"visibility,omitempty"`
	Type       string                 `json:"type,omitempty"` // type string,class,bool,etc.
}

func (p ValueProperties) Copy() ValueProperties {
	return ValueProperties{
		Visibility: p.Visibility,
		Type:       p.Type,
	}
}

// More detailed datatypes
type ValuePropVisibilityVar string

const (
	VisibilityPublic    ValuePropVisibilityVar = "public"
	VisibilityProtected ValuePropVisibilityVar = "protected"
	VisibilityPrivate   ValuePropVisibilityVar = "private"
	VisibilityUnknown   ValuePropVisibilityVar = ""
)
