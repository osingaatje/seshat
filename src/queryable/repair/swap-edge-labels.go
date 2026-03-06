package repair

import (
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/parse-result"
)

/*
	swap edge labels INSIDE THE STRUCT!

ex. of usefulness: if a student has dragged them to a location of where another label would be.

| class1 | * ------------- 1 | class2|

			(student drags * and 1 to other place:)
	 HOW IT VISUALLY   LOOKS: | class1 | 1 ------------- * | class2|
	 HOW IT INTERNALLY LOOKS: | class1 | * ------------- 1 | class2|

	it would be unfair to grade students based on code / data they cannot see. This is why we implement this.


	NOTE: This does not take into account dragging labels to other vertices
		ex.:  |class1| (A) ----------------- (B) |class2| (C) -------------- (D) |class3|
					(student drags B<->D and A<->C)
		HOW IT VISUALLY   LOOKS: |class1| (A) ----------------- (D) |class2| (C) -------------- (B) |class3|
		HOW IT INTERNALLY LOOKS: |class1| (C) ----------------- (B) |class2| (A) -------------- (D) |class3|

		we would consider the second case to be a pretty big 'skill issue'.
		If a student drags the labels within one edge because the software isn't cooperating, fine.
		..but if a student starts dragging labels to other vertices then we can consider that a failure in their thinking.
			(at least I (Douwe) think so. We should not account for such major failures)
*/
func swapEdgeLabels(c *context.Ctx, p *ParseResult) {
	if p == nil {
		c.LogWarn("Trying to swap edge labels on an empty parse result - BUG")
		return
	}

	c.LogErr("TODO MAKE SWAP LABELS")
}
