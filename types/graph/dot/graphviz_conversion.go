package displaygraph

import (
	. "github.com/osingaatje/seshat/types/graph/shared"
)

var EdgeStyleToGraphvizStyle map[EdgeLineStyle]string = map[EdgeLineStyle]string{
	EdgeLineStyleSolid:  "",
	EdgeLineStyleDotted: "dotted",
}

var ArrowStyleToGraphvizStyle map[ArrowStyleVariant]string = map[ArrowStyleVariant]string{
	ArrowStyleNoArrow:       "none",
	ArrowStyleArrow:         "normal",
	ArrowStyleHollowArrow:   "empty",
	ArrowStyleFilledDiamond: "diamond",
	ArrowStyleHollowDiamond: "odiamond",
}
