package helper

import (
	. "github.com/osingaatje/seshat/types/generic"
)

func GetCenterPos(p []Vector2D) Vector2D {
	// center between
	var centerPos Vector2D
	if len(p)%2 == 0 {
		// even
		halfPathLen := (len(p) + 1) / 2 // for 4 we want 2, but for 5 we want 3 (hence the +1)
		lMidPos, rMidPos := p[halfPathLen-1], p[halfPathLen]
		centerPos = lMidPos.Add(rMidPos.Sub(lMidPos).Div(2)) // A ---- B == A + (B-A)/2
	} else {
		// uneven
		centerPos = p[len(p)/2] // 3/2==1, 5/2==2 etc., correct for 0-indexed array
	}
	return centerPos
}
