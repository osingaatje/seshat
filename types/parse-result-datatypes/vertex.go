package parseresultdatatypes

import (
	"image/color"
	"reflect"
)

// VERTICES
type VertexProperty int

const (
	VertexPropClassType       VertexProperty = iota
	VertexPropClassVisibility                // "public", "private", "protected", ...
)

type VertexStyleProperty int

const (
	VertexStyleFillHex   VertexStyleProperty = iota // hex color!
	VertexStyleStrokeHex                            // hex color!
	VertexStyleStrokeWidth
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
