package repair

import (
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
	. "github.com/osingaatje/seshat/types/graph/shared"
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

const DISTANCE_PERC_THRESHOLD float64 = 0.35
const DIST_INF float64 = 999999999999

/*
 * Couple cases:
 * Start label swapped with either Middle or End label
 * Middle label swapped with Start/End label
 * End label swapped with Start/Middle label
 *
 * conditions:
 * - when the label is farther to one side of the edge than to its own side.
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

	labels := []**ParsedLabel{
		&e.FromProperties.Label, &e.Label, &e.ToProperties.Label,
	}
	labelTexts := []string{
		"From", "Middle", "To",
	}
	referencePosition := []Vector2D{
		fromNodePos, centerPos, toNodePos,
	}
	// we want to punish swapping the center node a bit because that's not often done and the distances are smaller.
	distance_scale := []float64{
		1, 1.5, 1,
	}
	nodeDist := fromNodePos.Dist(toNodePos)

	// close to start/middle/end
	distances := map[int][]float64{}
	for i, l := range labels {
		if *l == nil {
			distances[i] = []float64{DIST_INF, DIST_INF, DIST_INF}
			continue
		}
		distances[i] = []float64{}
		for j := range len(labels) {
			distances[i] = append(distances[i], labelDistance((*l).Location, referencePosition[j], distance_scale[j]))
		}
	}

	possibleSwaps := map[int][]int{
		0: {1, 2},
		1: {0, 2},
		2: {0, 1},
	}

	for nodeId := range len(labels) {
		swapIds := possibleSwaps[nodeId]

		_, canSwap := possibleSwaps[nodeId]
		if !canSwap || *(labels[nodeId]) == nil || distances[nodeId][nodeId] <= DISTANCE_PERC_THRESHOLD {
			continue // nodes that are nil or close or not swappable do not get swapped!
		}

		// if node is not close to its own spot,
		// try to swap with other index:
		var smallestDist float64 = DIST_INF
		var nodeToSwap int = -1
		for _, swapId := range swapIds {
			// choose smallest distance:
			_, isSwappable := possibleSwaps[swapId] // check if we can still swap the other node
			if isSwappable && distances[nodeId][swapId] < smallestDist &&
				distances[nodeId][swapId]/nodeDist <= DISTANCE_PERC_THRESHOLD &&
				distances[swapId][swapId]/nodeDist > DISTANCE_PERC_THRESHOLD {

				smallestDist, nodeToSwap = distances[nodeId][swapId], swapId
			}
		}
		if nodeToSwap == -1 {
			continue
		}

		swapLabels(
			c, labelTexts[nodeId], labelTexts[nodeToSwap],
			labels[nodeId], labels[nodeToSwap],
		)
		// don't forget to swap *back* the labels in our internal datastructure as well:
		labels[nodeId], labels[nodeToSwap] = labels[nodeToSwap], labels[nodeId]

		// mark as swapped:
		if labels[nodeId] != nil {
			delete(possibleSwaps, nodeId)
		}

		// we can still swap the other label if necessary!
		//if labels[nodeToSwap] != nil {
		//	delete(possibleSwaps, nodeToSwap)
		//}
	}
}

func labelDistance(lblPos Vector2D, referencePos Vector2D, scale float64) float64 { // , lengthOfEdge float64) float64 {
	// we detect a label swapped if its more than X% of the way to the other side
	return lblPos.Dist(referencePos) * scale
}

func swapLabels(c *context.Ctx, label1Txt string, label2Txt string, lbl1 **ParsedLabel, lbl2 **ParsedLabel) {
	c.LogDebug("Swapping labels %s and %s...", label1Txt, label2Txt)

	if *lbl1 != nil {
		if *lbl2 != nil {
			**lbl1, **lbl2 = **lbl2, **lbl1

		} else {
			tmpLbl1 := (*lbl1).Copy()
			*lbl1 = nil
			*lbl2 = tmpLbl1
		}
	} else /* lbl1 == nil */ {
		if *lbl2 != nil {
			*lbl1 = *lbl2
			*lbl2 = nil
		}
		// else case: both are nil, no swap needed.
	}
}
