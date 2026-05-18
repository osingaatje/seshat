package displaygraph

import (
	. "github.com/osingaatje/seshat/types/graph/shared"
)

var EdgeStyleToGraphvizStyle map[EdgeLineStyle]string = map[EdgeLineStyle]string{
	EdgeLineStyleSolid:  "",
	EdgeLineStyleDotted: "dotted",
}
