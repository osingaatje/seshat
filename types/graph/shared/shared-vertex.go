package shared

import (
	"image/color"

	. "github.com/osingaatje/seshat/types/generic"
)

// VERTICES
type VertexProperties struct {
	Type       string                 `json:"type,omitempty"`       // class, interface, ...
	Visibility ValuePropVisibilityVar `json:"visibility,omitempty"` // public, private, protected, ...
}

type VertexVisualProperties struct {
	// additional stuff for possible visualisation later on:
	Location               Vector2D   `json:"location"`
	Size                   Vector2D   `json:"size"`
	VertexStyleFillHex     color.RGBA `json:"fill_col"`
	VertexStyleStrokeHex   color.RGBA `json:"stroke_col"`
	VertexStyleStrokeWidth uint16     `json:"stroke_width"`
}
