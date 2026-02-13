package parseresultdatatypes

import (
	"image/color"
	"reflect"
)

// VERTICES
type VertexProperty string
type VertexStyleProperty string

const (
	VertexPropClassType       VertexProperty = "Class Type"
	VertexPropClassVisibility VertexProperty = "Visibility" // "public", "private", "protected", ...
	// TODO: add more if needed.

	//styling:
	VertexStyleFillHex     VertexStyleProperty = "Fill"   // hex color!
	VertexStyleStrokeHex   VertexStyleProperty = "Stroke" // hex color!
	VertexStyleStrokeWidth VertexStyleProperty = "Stroke Width"
)

// maps allowed properties to their data type
var VertexPropertyAll map[VertexProperty]reflect.Type = map[VertexProperty]reflect.Type{
	VertexPropClassType:       reflect.TypeOf(""),
	VertexPropClassVisibility: reflect.TypeOf(true),
}

var VertexStylePropertyAll map[VertexStyleProperty]reflect.Type = map[VertexStyleProperty]reflect.Type{
	VertexStyleFillHex:     reflect.TypeOf(color.RGBA{}),
	VertexStyleStrokeHex:   reflect.TypeOf(color.RGBA{}),
	VertexStyleStrokeWidth: reflect.TypeOf(uint16(0)),
}
