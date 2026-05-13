package repair

import (
	"math"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/generic"
	pr "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
	"github.com/osingaatje/seshat/types/repair"
)

// some edges may not be connected to both nodes. We hope that the creator of the diagram means to connect it to an edge, which we try to implement here.
func tryConnectEdgeEnds(c *context.Ctx,
	graph *pr.InternalGraph,
	iE *pr.InternalEdge) error {
	// if we've already added edge connections, we're done.
	if ((iE.FromId == nil) != (iE.FromEdgeId == nil)) && ((iE.ToId == nil) != (iE.ToEdgeId == nil)) {
		return nil
	}
	if len(iE.VisualProperties.Path) < 2 {
		panic("CANNOT HAVE LESS THAN TWO PATH POSITIONS")
	}

	if iE.FromId == nil && iE.FromEdgeId == nil {
		dist, newVId, newEId := getClosestVertexOrEdgeId(c, graph, iE, &iE.VisualProperties.Path[0])
		if newVId != INVALID_VERT_ID {
			iE.FromId = &newVId
		} else if newEId != INVALID_EDGE_ID {
			iE.FromEdgeId = &newEId
		} else if dist >= repair.DIST_INF {
			c.LogWarn("Could not find a close enough vertex/edge to connect to the start of edge '%d'", iE.Id)
		} else {
			panic("INVALID CODE: BUG - minimal distance calculation for connecting to a vertex/edge did produce a distance but not an ID")
		}
	}

	if iE.ToId == nil && iE.ToEdgeId == nil { // look for the last path location
		dist, newVId, newEId := getClosestVertexOrEdgeId(c, graph, iE, &iE.VisualProperties.Path[len(iE.VisualProperties.Path)-1])
		if newVId != INVALID_VERT_ID {
			iE.ToId = &newVId
		} else if newEId != INVALID_EDGE_ID {
			iE.ToEdgeId = &newEId
		} else if dist >= repair.DIST_INF {
			c.LogWarn("Could not find a close enough vertex/edge to connect to the end of edge '%d'", iE.Id)
		} else {
			panic("INVALID CODE: BUG - minimal distance calculation for connecting to a vertex/edge did produce a distance but not an ID")
		}
	}

	return nil
}

func getClosestVertexOrEdgeId(c *context.Ctx, graph *pr.InternalGraph, iE *pr.InternalEdge, location *Vector2D) (float64, VertexIdentifier, EdgeIdentifier) {
	var smallestDist float64 = repair.DIST_INF
	var smallestDistVertexId VertexIdentifier = INVALID_VERT_ID
	var smallestDistEdgeId EdgeIdentifier = INVALID_EDGE_ID

	for otherVertId, otherVertex := range graph.Vertices {
		if otherVertex.VisualProperties.Shape != VERTEX_SHAPE_RECT {
			c.LogWarn("Encountered a non-square vertex ('%d') while trying to repair edge '%d'", otherVertId, iE.Id)
			continue
		}
		// OLD WAY: calculate to distance from the center, then subtract half the size again (dirty hack)
		// centerNodePos := otherVertex.VisualProperties.Location.Add(otherVertex.VisualProperties.Size.Mul(0.5))
		// dist := centerNodePos.Dist(*location)
		// c.LogDebug("E'%d' dist. to V'%d' = %.2f", iE.Id, otherVertId, dist)

		dists, lengths := _getDistsForRectangularVertex(otherVertex, location)

		for i, d := range dists {
			l := lengths[i]

			normalisedDist := d / (l + 0.00001) /* to avoid NaN */ // previosuly: dist / ((w+h)/2)
			if normalisedDist < (float64(repair.VERTEX_CLOSENESS_PERCENT)/100) && normalisedDist < smallestDist {
				smallestDist = normalisedDist
				smallestDistVertexId = otherVertId
				smallestDistEdgeId = INVALID_EDGE_ID
			}
		}
	}

	for otherEdgeId, otherEdge := range graph.Edges {
		if iE.Id == otherEdgeId {
			continue
		}
		if len(otherEdge.VisualProperties.Path) < 2 {
			panic("CANNOT HAVE LESS THAN TWO PATH POSITIONS")
		}

		// make sublines of the path of the edge and try to see if we have a distance match:
		for i := range len(otherEdge.VisualProperties.Path) - 1 {
			startLineSegment, endLineSegment := otherEdge.VisualProperties.Path[i], otherEdge.VisualProperties.Path[i+1]

			dist, edgeLen := _calculateDistanceToLine(&startLineSegment, &endLineSegment, location)
			normalisedDist := dist / (edgeLen + 0.0001) /* avoid NaN */ // before: dist/edgelen
			if normalisedDist < (float64(repair.LINE_CLOSENESS_PERCENT)/100) && normalisedDist < smallestDist {
				smallestDist = normalisedDist
				smallestDistEdgeId = otherEdgeId
				smallestDistVertexId = INVALID_VERT_ID
			}
		}
	}

	return smallestDist, smallestDistVertexId, smallestDistEdgeId
}

func _getDistsForRectangularVertex(
	v *pr.InternalVertex,
	loc *Vector2D,
) ([]float64, []float64) {
	w, h := v.VisualProperties.Size.X, v.VisualProperties.Size.Y
	topleftX, topleftY := v.VisualProperties.Location.X, v.VisualProperties.Location.Y

	topLeft := Vector2D{X: topleftX, Y: topleftY}
	topRight := Vector2D{X: topleftX + w, Y: topleftY}
	bottRight := Vector2D{X: topleftX + w, Y: topleftY + h}
	bottLeft := Vector2D{X: topleftX, Y: topleftY + h}

	// if the location is inside of the vertex, it counts automatically.
	if loc.X > topLeft.X && loc.X < bottRight.X &&
		loc.Y > topLeft.Y && loc.Y < bottRight.Y {
		// location is inside of vertex. Return zero distance:
		return []float64{0}, []float64{0}
	}

	// else, calculate distances for each edge:

	dists := make([]float64, 4) // top, right, bottom, left
	lens := make([]float64, 4)

	dists[0], lens[0] /*topDist, topLen*/ = _calculateDistanceToLine(
		&topLeft, &topRight, loc,
	)
	dists[1], lens[1] /*rightDist, rightLen*/ = _calculateDistanceToLine(
		&topRight, &bottRight, loc,
	)
	dists[2], lens[2] /*bottDist, bottLen*/ = _calculateDistanceToLine(
		&bottLeft, &bottRight, loc,
	)
	dists[3], lens[3] /* leftDist, leftLen */ = _calculateDistanceToLine(
		&topLeft, &bottLeft, loc,
	)

	return dists, lens
}

// see https://mathworld.wolfram.com/Point-LineDistance2-Dimensional.html
// For line segments: https://www.geeksforgeeks.org/dsa/minimum-distance-from-a-point-to-the-line-segment-using-vectors/
func _calculateDistanceToLine( /*start of line*/ a *Vector2D /*end of line*/, b *Vector2D, loc *Vector2D) (dist float64, edgeLength float64) {
	if a == nil || b == nil || loc == nil {
		panic("BUG IN CODE - FIX YOUR SHIT")
	}

	// sqrt(dx^2 + dy^2)
	lenOfEdge := math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2)) // <>---------------<> the length inbetween lineStart and lineEnd

	ba := b.Sub(*a) // (end.X - start.X, end.Y - start.Y)

	aLoc := loc.Sub(*a)
	bLoc := loc.Sub(*b)

	// calculate dot products to see if the point 'loc' lies on the projection of the line.
	// if this is the case, we return the distance to the line.
	// if this is *not* the case, we return the minimum distance to each individual point.

	// dot products
	ABDotBLoc := ba.X*bLoc.X + ba.Y*bLoc.Y
	ABDotALoc := ba.X*aLoc.X + ba.Y*aLoc.Y

	if ABDotBLoc > 0 {
		return math.Sqrt(math.Pow(loc.X-b.X, 2) + math.Pow(loc.Y-b.Y, 2)), lenOfEdge
	} else if ABDotALoc < 0 {
		return math.Sqrt(math.Pow(loc.X-a.X, 2) + math.Pow(loc.Y-a.Y, 2)), lenOfEdge
	} else if lenOfEdge == 0 {
		return 0, 0
	} else {
		return math.Abs(ba.X*aLoc.Y-ba.Y*aLoc.X) / lenOfEdge, lenOfEdge
	}
}
