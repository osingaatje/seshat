package parseresultdatatypes

import (
	"image/color"
)

// VERTICES
type VertexProperties struct {
	Type       string                 `json:"type,omitempty"`       // class, interface, ...
	Visibility ValuePropVisibilityVar `json:"visibility,omitempty"` // public, private, protected, ...
}

type VertexStyleProperties struct {
	VertexStyleFillHex     color.RGBA `json:"fill_col"`
	VertexStyleStrokeHex   color.RGBA `json:"stroke_col"`
	VertexStyleStrokeWidth uint16     `json:"stroke_width"`
}
