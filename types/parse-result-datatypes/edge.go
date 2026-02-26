package parseresultdatatypes

import (
	"github.com/osingaatje/seshat/helper"
)

// EDGES
type EdgeEndProperties struct { // the properties for the end of an edge ( <>------> )
	ArrowStyle   ArrowStyleVariant    `json:"arrowstyle,omitempty"`
	Multiplicity *helper.Multiplicity `json:"multiplicity,omitempty"`
}

type EdgeStyleProperties struct {
	LineStyle EdgeLineStyle `json:"line_style,omitempty"`
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
