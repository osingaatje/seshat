package shared

import (
	"github.com/osingaatje/seshat/types/generic"
)

type VertexIdentifier uint32
type EdgeIdentifier uint32

const INVALID_VERT_ID = VertexIdentifier(1<<32 - 1)
const INVALID_EDGE_ID = EdgeIdentifier(1<<32 - 1)

func NewVertexIdentifierInt(nId *int) *VertexIdentifier {
	if nId == nil {
		return nil
	}
	res := VertexIdentifier(uint32(*nId))
	return &res
}
func NewVertexIdentifier(nId *int) *VertexIdentifier {
	if nId == nil {
		return nil
	}
	id := VertexIdentifier(uint32(*nId))
	return &id
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
type Label struct {
	Text     string           `json:"text"`
	Location generic.Vector2D `json:"location"`
}

func (l *Label) Copy() *Label {
	if l == nil {
		return nil
	}
	return &Label{
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
