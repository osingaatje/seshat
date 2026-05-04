package repair

import (
	"fmt"
	"math"

	. "github.com/osingaatje/seshat/types/generic"
	pr "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
	"github.com/osingaatje/seshat/types/repair"
)

// some edges may not be connected to both nodes. We hope that the creator of the diagram means to connect it to an edge, which we try to implement here.
func tryConnectEdgeEnds(graph *pr.InternalGraph, iE *pr.InternalEdge) error {
	// if we've already added edge connections, we're done.
	if ((iE.FromId == nil) != (iE.FromEdgeId == nil)) && ((iE.ToId == nil) != (iE.ToEdgeId == nil)) {
		return nil
	}
	if len(iE.VisualProperties.Path) < 2 {
		panic("CANNOT HAVE LESS THAN TWO PATH POSITIONS")
	}

	// just so that we don't have to make another function.
	for range 2 {
		var vertexIdToFix **VertexIdentifier
		var edgeIdToFix **EdgeIdentifier
		var location *Vector2D
		var edgeName string

		if iE.FromId == nil && iE.FromEdgeId == nil {
			vertexIdToFix = &iE.FromId
			edgeIdToFix = &iE.FromEdgeId
			location = &iE.VisualProperties.Path[0]
			edgeName = "start"
		} else if iE.ToId == nil && iE.ToEdgeId == nil {
			vertexIdToFix = &iE.ToId
			edgeIdToFix = &iE.ToEdgeId
			location = &iE.VisualProperties.Path[len(iE.VisualProperties.Path)-1]
			edgeName = "end"
		} else {
			break
		}

		var smallestDist float64 = math.Inf(1)
		var smallestDistVertexId VertexIdentifier = INVALID_VERT_ID
		var smallestDistEdgeId EdgeIdentifier = INVALID_EDGE_ID

		for otherVertId, otherVertex := range graph.Vertices {
			h, w := otherVertex.VisualProperties.Size.X, otherVertex.VisualProperties.Size.Y

			// calculate to distance from the center, then subtract half the size again (dirty hack)
			topLeftPos := otherVertex.VisualProperties.Location
			centerNodePos := topLeftPos.Add(otherVertex.VisualProperties.Size.Mul(0.5))

			dist := centerNodePos.Dist(*location)

			normalisedDist := dist / ((h + w) / 2)
			if normalisedDist < float64(repair.VERTEX_CLOSENESS_PERCENT)/100 && normalisedDist < smallestDist {
				smallestDist = normalisedDist
				smallestDistVertexId = otherVertId
				smallestDistEdgeId = INVALID_EDGE_ID
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
				normalisedDist := dist / edgeLen
				if normalisedDist < float64(repair.LINE_CLOSENESS_PERCENT)/100 && normalisedDist < smallestDist {
					smallestDist = normalisedDist
					smallestDistEdgeId = otherEdgeId
					smallestDistVertexId = INVALID_VERT_ID
				}
			}
		}

		if smallestDist < 99999999 {
			if smallestDistVertexId != INVALID_VERT_ID {
				*vertexIdToFix = &smallestDistVertexId
			} else if smallestDistEdgeId != INVALID_EDGE_ID {
				*edgeIdToFix = &smallestDistEdgeId
			} else {
				panic("INVALID CODE: BUG - minimal distance calculation for connecting to a vertex/edge did produce a distance but not an ID")
			}
		} else {
			return fmt.Errorf("Could not find a close enough edge to connect %s point of edge '%d' to located at (%.2f,%.2f)", edgeName, iE.Id, location.X, location.Y)
		}
	}
	return nil
}

// see https://mathworld.wolfram.com/Point-LineDistance2-Dimensional.html
func _calculateDistanceToLine(lineStart *Vector2D, lineEnd *Vector2D, loc *Vector2D) (dist float64, edgeLength float64) {
	// calculate distance between the point 'loc(ation)' = (x_0, y_0) and the line lineStart<->lineEnd = (x_1,y_1), (x_2, y_2)
	// we use this formula: d = | v^ . r | = ( |(x_2 - x_1)*(y_1 - y_0) - (x_1 - x_0)(y_2 - y_1)| ) \ ( \sqrt( (x_2-x_1)^2 + (y_2-y_1)^2 ) )

	x_0, y_0 := loc.X, loc.Y
	x_1, y_1 := lineStart.X, lineStart.Y
	x_2, y_2 := lineEnd.X, lineEnd.Y

	lenOfEdge := math.Sqrt(math.Pow(x_2-x_1, 2) + math.Pow(y_2-y_1, 2)) // <>---------------<> the length inbetween lineStart and lineEnd

	topPart := math.Abs((x_2-x_1)*(y_1-y_0) - (x_1-x_0)*(y_2-y_1))
	return topPart / lenOfEdge, lenOfEdge
}
