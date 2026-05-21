package shared

import (
	"fmt"

	"github.com/osingaatje/seshat/helper/multiplicity"
	"github.com/osingaatje/seshat/types/generic"
)

// EDGES
type EdgeEndProperties struct { // the properties for the end of an edge ( <>------> )
	ArrowStyle ArrowStyleVariant `json:"arrowstyle"`
	Label      *Label            `json:"label,omitempty"` // don't convert to stricter representation (multiplicity or something) yet, as we want to repair swapped labels etc.
	// multiplicity conversion should happen later! Multiplicity *helper.Multiplicity `json:"multiplicity,omitempty"`
}

type InternalEdgeEndProperties struct {
	ArrowStyle   ArrowStyleVariant          `json:"arrowstyle"`
	Multiplicity *multiplicity.Multiplicity `json:"multiplicity,omitempty"`
}

func (e EdgeEndProperties) ToInternal() (props InternalEdgeEndProperties, err error) {
	props = InternalEdgeEndProperties{
		ArrowStyle:   e.ArrowStyle,
		Multiplicity: nil,
	}
	if e.Label.HasText() {
		mult, ok := multiplicity.GetMultiplicity(e.Label.Text)
		if !ok {
			return props, fmt.Errorf("Could not parse multiplicity '%s'", e.Label.Text)
		}
		props.Multiplicity = &mult
	}
	return props, nil
}

func (p EdgeEndProperties) Copy() EdgeEndProperties {
	res := EdgeEndProperties{
		ArrowStyle: p.ArrowStyle,
		Label:      p.Label.Copy(), // since it's a pointer we need to carefully copy it
	}
	return res
}

type EdgeStyleProperties struct {
	LineStyle EdgeLineStyle `json:"line_style"`
}

func (s EdgeStyleProperties) Copy() EdgeStyleProperties {
	return EdgeStyleProperties{
		LineStyle: s.LineStyle,
	}
}

type EdgeVisualProperties struct {
	Path []generic.Vector2D `json:"path"`
}

func (s EdgeVisualProperties) Copy() EdgeVisualProperties {
	res := EdgeVisualProperties{
		Path: make([]generic.Vector2D, len(s.Path)),
	}
	copy(res.Path, s.Path)
	return res
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

func ArrowStyleIsDirected(a ArrowStyleVariant) bool {
	switch a {
	case ArrowStyleNoArrow:
		return false
	case ArrowStyleArrow:
		return true
	case ArrowStyleHollowArrow:
		return true
	case ArrowStyleFilledDiamond:
		return true
	case ArrowStyleHollowDiamond:
		return true
		// add new arrow styles here
	}
	panic("Unknown arrow style!!! Bug in code fix pls")
}

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
