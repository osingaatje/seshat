package parseresult

import (
	utml "github.com/osingaatje/seshat/types/parse-result-utml"
)

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
	VisibilityUnknown   ValuePropVisibilityVar = "unknown"
)

// location on some grid or whatever
type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2D) New(utmlPos utml.UTMLXY) Vector2D {
	return Vector2D{
		X: float64(utmlPos.X),
		Y: float64(utmlPos.Y),
	}
}
func (v Vector2D) NewInt(x int, y int) Vector2D {
	return Vector2D{
		X: float64(x),
		Y: float64(y),
	}
}
func (v Vector2D) Add(vO Vector2D) Vector2D {
	return Vector2D{
		X: v.X + vO.X,
		Y: v.Y + vO.Y,
	}
}
func (v Vector2D) Div(factor float64) Vector2D {
	return Vector2D{
		X: v.X / factor,
		Y: v.Y / factor,
	}
}
func (v Vector2D) Mul(factor float64) Vector2D {
	return Vector2D{
		X: v.X * factor,
		Y: v.Y * factor,
	}
}
