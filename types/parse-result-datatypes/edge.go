package parseresultdatatypes

import (
	"reflect"
)

// EDGES
type EdgeEndProperty int // the properties for the end of an edge ( <>------> )

const (
	EdgeEndPropArrowStyle EdgeEndProperty = iota
	EdgeEndPropMultiplicity
)

var EdgeEndPropertyAll map[EdgeEndProperty]reflect.Type = map[EdgeEndProperty]reflect.Type{
	EdgeEndPropArrowStyle:   reflect.TypeOf(ArrowStyleNoArrow),
	EdgeEndPropMultiplicity: reflect.TypeOf(""),
}

type EdgeStyleProperty int // styling specifics (for visualisation)
const (
	// styling
	EdgeStyleLine EdgeStyleProperty = iota
)

var EdgeStylePropertyAll map[EdgeStyleProperty]reflect.Type = map[EdgeStyleProperty]reflect.Type{
	EdgeStyleLine: reflect.TypeOf(EdgeLineStyleSolid),
}

// more detailed datatypes
type ArrowStyleVariant int

const (
	ArrowStyleNoArrow ArrowStyleVariant = iota
	ArrowStyleArrow
	ArrowStyleHollowArrow
	ArrowStyleFilledDiamond
	ArrowStyleHollowDiamond
)

type EdgeLineStyle int

const (
	EdgeLineStyleSolid EdgeLineStyle = iota
	EdgeLineStyleDotted
)
