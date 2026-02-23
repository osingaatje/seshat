package parse

import (
	. "github.com/osingaatje/seshat/types"
	. "github.com/osingaatje/seshat/types/parse-result-datatypes"
)

var UTMLLineStyleToParsedStyle map[UTMLLineStyle]EdgeLineStyle = map[UTMLLineStyle]EdgeLineStyle{
	UTMLLineStyleDashed: EdgeLineStyleDotted,
	UTMLLineStyleDotted: EdgeLineStyleDotted,
	UTMLLineStyleFilled: EdgeLineStyleSolid,
}
