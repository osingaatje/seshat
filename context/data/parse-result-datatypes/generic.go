package parseresultdatatypes

import (
	"reflect"
)

// Properties of values, e.g.: "field" : Visibility=private, Type=string
type ValueProperty string

const (
	ValuePropVisibility ValueProperty = "Visibility"
	ValuePropType       ValueProperty = "Type"
)

var ValuePropertyAll map[ValueProperty]reflect.Type = map[ValueProperty]reflect.Type{
	ValuePropVisibility: reflect.TypeOf(VisibilityPublic),
	ValuePropType:       reflect.TypeOf(""), // we allow anything here, class, inherited, ...
}

// More detailed datatypes
type ValuePropVisibilityVar string

const (
	VisibilityPublic    ValuePropVisibilityVar = "Public"
	VisibilityProtected ValuePropVisibilityVar = "Protected"
	VisibilityPrivate   ValuePropVisibilityVar = "Private"
	VisibilityUnknown   ValuePropVisibilityVar = "?"
)
