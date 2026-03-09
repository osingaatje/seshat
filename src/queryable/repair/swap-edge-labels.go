package repair

import (
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/generic"
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

	for _, e := range p.Edges {
		swapEdgeLabelsForEdge(c, p, e)
	}
}

/*
 * Couple cases:
 * Start label swapped with either Middle or End label
 * Middle label swapped with Start/End label
 * End label swapped with Start/Middle label
 *
 * conditions:
 * - when the label is farther to one side of the edge than to its own side.
 *
 * NOTE: We limit ourselves to one swap.
 */
func swapEdgeLabelsForEdge(c *context.Ctx, p *ParseResult, e *ParsedEdge) {
	if e == nil {
		return
	}

	fromNode, okFrom := p.Vertices[e.FromId]
	toNode, okTo := p.Vertices[e.ToId]
	if !okFrom || !okTo {
		c.LogErr("Invalid Parse Result: missing vertices (trying to swap edge labels)")
		return
	}

	fromNodePos := fromNode.VisualProperties.Location.Add(fromNode.VisualProperties.Size.Div(2)) // center of the node
	toNodePos := toNode.VisualProperties.Location.Add(toNode.VisualProperties.Size.Div(2))       // center of the node

	// center between A ---- B == A + (B-A)/2
	centerPos := fromNodePos.Add(toNodePos.Sub(fromNodePos).Div(2))
	// distance between from and to node
	dist := fromNodePos.Dist(toNodePos)

	fromLbl := e.FromProperties.Label
	midLbl := e.Label
	toLbl := e.ToProperties.Label

	fromLblFarAway := fromLbl != nil && labelTooFarAway(
		fromLbl.Location,
		fromNodePos,
		dist)

	midLblFarAway := midLbl != nil && labelTooFarAway(
		midLbl.Location,
		centerPos,
		dist/2) // (middle label should be at the center, so the distance should be halved to reach either node)

	toLblFarAway := toLbl != nil && labelTooFarAway(
		toLbl.Location,
		toNodePos,
		dist)

	// from lbl <--> [mid, to]
	if fromLblFarAway {
		if toLblFarAway {
			swapLabels(c, "From", "To", &fromLbl, &toLbl)
		} else if midLblFarAway {
			swapLabels(c, "From", "Middle", &fromLbl, &midLbl)
		}

		// mid lbl <--> [to, from]
	} else if midLblFarAway {
		if fromLblFarAway {
			swapLabels(c, "Middle", "From", &midLbl, &fromLbl)
		} else if toLblFarAway {
			swapLabels(c, "Middle", "To", &midLbl, &toLbl)
		}

		// to lbl <--> [from, mid]
	} else if toLblFarAway {
		if fromLblFarAway {
			swapLabels(c, "To", "From", &toLbl, &fromLbl)
		} else if midLblFarAway {
			swapLabels(c, "To", "Middle", &toLbl, &midLbl)
		}
	}

}

func labelTooFarAway(lblPos Vector2D, referencePos Vector2D, lengthOfEdge float64) bool {
	// we detect a label swapped if its more than 75% of the way to the other side
	// .. but somehow the offset position is not exactly the same factor, so we'll just do somewhere around 20-40% I guess.
	return lblPos.Dist(referencePos)/lengthOfEdge > 0.25
}

func swapLabels(c *context.Ctx, label1Txt string, label2Txt string, lbl1 **ParsedLabel, lbl2 **ParsedLabel) {
	c.LogDebug("Swapping labels %s and %s...", label1Txt, label2Txt)
	*lbl1, *lbl2 = *lbl2, *lbl1
}
