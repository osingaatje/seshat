package convert

import (
	. "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/parse-result"
	. "github.com/osingaatje/seshat/types/parse-result-utml"
)

var UTMLClassTypeToInternal map[UTMLClassType]string = map[UTMLClassType]string{
	UTMLClassTypeClass:     "class",
	UTMLClassTypeAbstract:  "abstract",
	UTMLClassTypeInterface: "interface",
}

var UTMLLineStyleToParsedStyle map[UTMLLineStyle]EdgeLineStyle = map[UTMLLineStyle]EdgeLineStyle{
	UTMLLineStyleDashed: EdgeLineStyleDotted,
	UTMLLineStyleDotted: EdgeLineStyleDotted,
	UTMLLineStyleFilled: EdgeLineStyleSolid,
}

var UTMLVisibilityToInternalVisibility map[UTMLVisibility]ValuePropVisibilityVar = map[UTMLVisibility]ValuePropVisibilityVar{
	UTMLVisibilityPrivate:   VisibilityPrivate,
	UTMLVisibilityProtected: VisibilityProtected,
	UTMLVisibilityPackage:   VisibilityProtected, // Package -> Protected
	UTMLVisibilityPublic:    VisibilityPublic,
}

var UTMLArrowStyleToInteral map[UTMLArrowHeadStyle]ArrowStyleVariant = map[UTMLArrowHeadStyle]ArrowStyleVariant{
	UTMLArrowStyleNone:               ArrowStyleNoArrow,
	UTMLArrowStyleSmallFilledArrow:   ArrowStyleArrow,
	UTMLArrowStyleFilledDiamond:      ArrowStyleFilledDiamond,
	UTMLArrowStyleUnfilledDiamond:    ArrowStyleHollowDiamond,
	UTMLArrowStyleLargeUnfilledArrow: ArrowStyleHollowArrow,
}
