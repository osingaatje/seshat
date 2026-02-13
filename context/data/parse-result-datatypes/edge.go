package parseresultdatatypes

import (
	"reflect"
)

// EDGES
type EdgeEndProperty string // the properties for the end of an edge ( <>------> )

const (
	EdgeEndPropArrowStyle   EdgeEndProperty = "Arrow Style"
	EdgeEndPropMultiplicity EdgeEndProperty = "Multiplicity"
)

var EdgeEndPropertyAll map[EdgeEndProperty]reflect.Type = map[EdgeEndProperty]reflect.Type{
	EdgeEndPropArrowStyle:   reflect.TypeOf(ArrowStyleNoArrow),
	EdgeEndPropMultiplicity: reflect.TypeOf(""),
}

type EdgeProperty string      // general properties of the edge
type EdgeStyleProperty string // styling specifics (for visualisation)
const (
	// core properties
	EdgePropLabelText EdgeProperty = "Label"

	// styling
	EdgeStyleLine EdgeStyleProperty = "LineStyle"
)

var EdgePropertyAll map[EdgeProperty]reflect.Type = map[EdgeProperty]reflect.Type{
	EdgePropLabelText: reflect.TypeOf(""),
}
var EdgeStylePropertyAll map[EdgeStyleProperty]reflect.Type = map[EdgeStyleProperty]reflect.Type{
	EdgeStyleLine: reflect.TypeOf(EdgeLineStyleSolid),
}

// more detailed datatypes
type ArrowStyleVariant string

const (
	ArrowStyleNoArrow       ArrowStyleVariant = "None"
	ArrowStyleArrow         ArrowStyleVariant = "Arrow"
	ArrowStyleHollowArrow   ArrowStyleVariant = "HollowArrow"
	ArrowStyleFilledDiamond ArrowStyleVariant = "Diamond"
	ArrowStyleHollowDiamond ArrowStyleVariant = "HollowDiamond"
)

type EdgeLineStyle string

const (
	EdgeLineStyleSolid  EdgeLineStyle = "Solid"
	EdgeLineStyleDotted EdgeLineStyle = "Dotted"
)
