package parseresult

// EDGES
type EdgeEndProperties struct { // the properties for the end of an edge ( <>------> )
	ArrowStyle ArrowStyleVariant `json:"arrowstyle"`
	Label      *ParsedLabel      `json:"label,omitempty"` // don't convert to stricter representation (multiplicity or something) yet, as we want to repair swapped labels etc.
	// multiplicity conversion should happen later! Multiplicity *helper.Multiplicity `json:"multiplicity,omitempty"`
}

type EdgeStyleProperties struct {
	LineStyle EdgeLineStyle `json:"line_style"`
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
