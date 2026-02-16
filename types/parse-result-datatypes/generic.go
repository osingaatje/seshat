package parseresultdatatypes

import (
	"reflect"
)

// Properties of values, e.g.: "field" : Visibility=private, Type=string
type ValueProperty int

const (
	ValuePropVisibility ValueProperty = iota
	ValuePropType
)

var ValuePropertyAll map[ValueProperty]reflect.Type = map[ValueProperty]reflect.Type{
	ValuePropVisibility: reflect.TypeOf(VisibilityPublic),
	ValuePropType:       reflect.TypeOf(""), // we allow anything here, class, inherited, ...
}

// More detailed datatypes
type ValuePropVisibilityVar int

const (
	VisibilityPublic ValuePropVisibilityVar = iota
	VisibilityProtected
	VisibilityPrivate
	VisibilityUnknown
)
