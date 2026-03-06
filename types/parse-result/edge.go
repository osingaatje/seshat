package parseresult

import (
	. "github.com/osingaatje/seshat/types/generic"
)

// EDGES
type EdgeEndProperties struct { // the properties for the end of an edge ( <>------> )
	ArrowStyle ArrowStyleVariant `json:"arrowstyle"`
	Label      *ParsedLabel      `json:"label,omitempty"` // don't convert to stricter representation (multiplicity or something) yet, as we want to repair swapped labels etc.
	// multiplicity conversion should happen later! Multiplicity *helper.Multiplicity `json:"multiplicity,omitempty"`
}

func (p EdgeEndProperties) Copy() EdgeEndProperties {
	res := EdgeEndProperties{
		ArrowStyle: p.ArrowStyle,
		Label:      nil, // since it's a pointer we need to carefully copy it
	}
	if p.Label != nil {
		res.Label = p.Label.Copy()
	}
	return res
}

type EdgeStyleProperties struct {
	LineStyle     EdgeLineStyle `json:"line_style"`
	StartLocation Vector2D      `json:"start_location"`
	EndLocation   Vector2D      `json:"end_location"`
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

type EdgeLabelPos int

const (
	EdgeLabelPosStart EdgeLabelPos = iota
	EdgeLabelPosMiddle
	EdgeLabelPosEnd
)
