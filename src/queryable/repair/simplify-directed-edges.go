package repair

import (
	"fmt"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

/*
 * This function aims to make individual edges of weirdly combined edges (perhaps by people that are less familiar with modelling tools)
 *
 * We want to correct this type of graph: (edge ends are indicated with '*')
 * 						[V1]
 * 						  ^
 * 						  |
 * 				[V3]------*---[V2]
 *                        \
 *                         --*----[V4]
 * ... and introduce individual edges for this, since we *ASSUME* that the diagram creator intends to view it as multiple individual edges connected from [V2],[V3],[V4] to [V1]
 *
 * 				  /-->[V1]
 * 				[V3]   ^ ^
 *                     |  \__[V2]
 *                     |
 *                      \___[V4]
 *  (note that because of ASCII limitations, we have changed the path locations. Ideally, the simplification preserves the locations)
 */
func simplifyDirectedEdges(c *context.Ctx, g *InternalGraph) error {
	// Filter on the directed edge ends that could start a combined edge
	directedEdgeEnds := helper.FilterMap(g.Edges, func(id EdgeIdentifier, e *InternalEdge) bool {
		// [V]<---[E] or [E]--->[V]
		return (ArrowStyleIsDirected(e.FromProperties.ArrowStyle) && e.FromId != nil && e.FromEdgeId == nil && e.ToId == nil && e.ToEdgeId != nil) ||
			(e.FromId == nil && e.FromEdgeId != nil && ArrowStyleIsDirected(e.ToProperties.ArrowStyle) && e.ToId != nil && e.ToEdgeId == nil)
	})

	// nothing to do
	if len(directedEdgeEnds) == 0 {
		return nil
	}

	// per directed edge end we can have new simplified edge parts
	newEdgeParts := map[EdgeIdentifier][][]EdgeIdentifier{}

	// explore the other edges
	for eId, _ := range directedEdgeEnds {
		// list of edges
		newSimplifiedEdgeParts := [][]EdgeIdentifier{
			{eId},
		}
		progress := []bool{true}

		for {
			if !helper.AnyBool(progress) {
				break
			}

			for edgePartsIndex, edgeIds := range newSimplifiedEdgeParts {
				if !progress[edgePartsIndex] {
					continue
				}

				lastEId := edgeIds[len(edgeIds)-1]
				lineStyle := g.Edges[lastEId].StyleProperties.LineStyle

				// find all edge-to-edge connections
				discoveredVertexIds, discoveredEdgeIds := g.ConnectionsForEdge(lastEId)

				// filter out the existing edges in the path, as well as edges that have a direction arrow
				skipMessages := []string{}
				alreadyDiscoveredEdgeIndices := []int{}
				for i, e := range discoveredEdgeIds {
					if slices.Contains(edgeIds, e) {
						alreadyDiscoveredEdgeIndices = append(alreadyDiscoveredEdgeIndices, i)
						continue
					}

					// SANITY CHECKS
					edge := g.Edges[e]
					if ArrowStyleIsDirected(edge.FromProperties.ArrowStyle) || ArrowStyleIsDirected(edge.ToProperties.ArrowStyle) {
						skipMessages = append(skipMessages, fmt.Sprintf("An arrow of edge '%d' was directed, so we cannot combine this entire series of edges", e))
						break
					}
					if edge.StyleProperties.LineStyle != lineStyle {
						skipMessages = append(skipMessages, fmt.Sprintf("Edge '%d' has a different line style, so we cannot combine this entire series of edges", e))
						break
					}
				}
				if len(skipMessages) > 0 {
					skipStr := strings.Join(skipMessages, ",")
					// progress[edgePartsIndex] = false
					return fmt.Errorf("Cannot complete simplification of edges because of edges starting at edge ID '%d': %s", edgeIds[0], skipStr)
				}

				// delete already discovered edges
				for _, ind := range alreadyDiscoveredEdgeIndices {
					discoveredEdgeIds = slices.Delete(discoveredEdgeIds, ind, ind+1)
				}

				if len(discoveredVertexIds) > 1 {
					panic("I thought we were looking for edge-to-vertex connections? Something is wrong with connectionsToEdge")
				}
				if len(discoveredVertexIds) > 0 { // if we've hit a vertex after the first edge, we know that this edge cannot be explored anymore.
					progress[edgePartsIndex] = false
				}

				// for each edge connection, add a new possible edge
				for edgeIndex, e := range discoveredEdgeIds {
					edge := g.Edges[e]
					if edgeIndex == 0 {
						// we cannot combine an edge connected to a vertex, with a label on the edge-connected side (i.e. [V]----<lbl>--*----....), but otherwise ([V]<lbl>----*---...) it is allowed
						if edge.Label != nil || (edge.FromId != nil && edge.ToProperties.Label != nil) ||
							(edge.ToId != nil && edge.FromProperties.Label != nil) {
							// progress[edgePartsIndex] = false
							return fmt.Errorf("Cannot simplify edges: edge '%d' has a start/middle/end label", e)
						}

						// for the first explored connection, add it to the existing arrow.
						newSimplifiedEdgeParts[edgeIndex] = append(newSimplifiedEdgeParts[edgeIndex], e)
						progress[edgePartsIndex] = true

					} else {
						if edge.FromProperties.Label != nil || edge.Label != nil || edge.ToProperties.Label != nil {
							// progress[edgePartsIndex] = false
							// skip = true
							return fmt.Errorf("Cannot simplify edges: Edge '%d' has a start/middle/end label", e)
						}

						// for all extra edges (branching), add new progresses etc.
						newEdge := append(slices.Clone(edgeIds), e)

						// add new progressing edge
						newSimplifiedEdgeParts = append(newSimplifiedEdgeParts, newEdge)
						progress = append(progress, true)
					}
				}
			}
		}

		/* we only collect the edge parts here, because there could be a situation where two directed edges go to the same 'network' of other edges, i.e.:
				[V1]  [V2]
			 	 ^		^
				 |		|
				 \--*--/
					|
		     [V3]___/\___[V4]

			..in this case, we would destroy our algorithm if we already deleted the old edges and connected the new ones, because either only V1 or V2 would have correct edges pointing to it.
			Instead, we only register the new edge parts, and thus finally arrive at the edge sets [V3]->[V1],[V4]->[V1] AND [V3]->[V2],[V4]->[V2]
		*/
		if len(newSimplifiedEdgeParts) > 0 {
			newEdgeParts[eId] = newSimplifiedEdgeParts
		}
	}

	// nothing to do
	if len(newEdgeParts) == 0 {
		return nil
	}

	// consolidate the new edges

	// keep the old edge info so that we can safely delete the edge entries from the graph itself without worrying about losing track of the edge information for other edges to be consolidated
	oldEdges := map[EdgeIdentifier]bool{}
	for _, edgeIdArrays := range newEdgeParts {
		for _, eIds := range edgeIdArrays {
			// keep old information on the edges
			for _, eId := range eIds {
				oldEdges[eId] = true
			}

			firstEdge := g.Edges[eIds[0]]
			firstEdgeFlipped := false
			firstId := firstEdge.FromId
			if firstId == nil {
				firstId = firstEdge.ToId
				firstEdgeFlipped = true
			}
			if firstId == nil {
				panic("Cannot have an edge with only edge connections!! BUG in algorithm")
			}

			lastEdge := g.Edges[eIds[len(eIds)-1]]
			lastId := lastEdge.FromId
			if lastId == nil {
				lastId = lastEdge.ToId
			}
			if lastId == nil {
				panic("Cannot have an edge with only edge connections!! BUG in algorithm")
			}

			combinedEdge := InternalEdge{
				Id:             g.NewEdgeId(),
				FromId:         firstId,
				ToId:           lastId,
				FromProperties: firstEdge.FromProperties.Copy(),
				ToProperties:   lastEdge.ToProperties.Copy(),
				// explicitly no label
				Label:           nil,
				StyleProperties: lastEdge.StyleProperties,
				VisualProperties: EdgeVisualProperties{
					Path: nil,
				},
			}

			newPath := slices.Clone(firstEdge.VisualProperties.Path)
			if firstEdgeFlipped {
				slices.Reverse(newPath)
			}

			// make the combined path by going one by one FROM first edge 0 TO last edge
			for i, eId := range eIds {
				if i == 0 {
					continue // was already added in initialisation
				}
				e := g.Edges[eId]
				path := slices.Clone(e.VisualProperties.Path)
				// check if we need to flip the path (edges may be constructed in both directions, so ---e[i+1-->e[i] or <---e[i+1]----e[i])
				if e.ToEdgeId != nil && (*e.ToEdgeId) == eIds[i-1] {
					slices.Reverse(path)
				}

				for i, vec := range path {
					//skip starting vector because we don't want a zig-zag path
					if i == 0 {
						continue
					}
					if newPath[len(newPath)-1] == vec {
						continue
					}
					newPath = append(newPath, vec)
				}
			}
			combinedEdge.VisualProperties.Path = newPath

			// add it to the graph:
			g.Edges[combinedEdge.Id] = &combinedEdge
		}
	}

	// REMOVE OLD EDGES
	for eId, _ := range oldEdges {
		delete(g.Edges, eId)
	}

	return nil
}
