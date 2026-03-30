package grade

import (
	"math"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/internal-rep"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

/*
 * Idea (inspired by Smith (2004), Thomas (2009), Vachharajani (2014), Bian (2020))
 * (steps inspired by Smith 2004, skipping segmentation and assimilation because the internal graph repr. already does that)
 *
 * 1. Identify / match: using (Minimal) Meaningful Units, with algorithms such as levenshtein distance (syntactic) and WordNet similarity score (HSO/WUP/LIN) (semantic)
 * 2. Aggregate: 		continually increase Meaningful Unit until it can no longer be aggregated) using structural graph matching
 * 3. Interpret:        produce grades based on present/missing/extra elements and parts of those elements
 */
func gradeDiag(c *context.Ctx, cmd GradeCmd) *GradeResult {
	if cmd.ReferenceSolution == nil || cmd.Submission == nil || cmd.Rubric == nil {
		c.LogErr("Referencesolution, submission, or rubric was not present when grading diagram!")
		return nil
	}

	potentialMatches, syntacticDistance, semanticMiniLM := getPotentialMatchingVertices(c, cmd)

	/*
	 * This map contains the minimal meaningful units, mapped from REFERENCE -> SUBMISSION
	 */
	minimalUnits := map[*InternalGraph]*InternalGraph{}
	fixedIds := map[ /*reference*/ VertexIdentifier] /*submission*/ VertexIdentifier{}

	for refId, subIdMap := range potentialMatches {
		var smallestSyntacticDistance int = math.MaxInt
		var highestSemanticMiniScore float64 = math.MaxFloat64
		var highestScoreId VertexIdentifier = 0

		for subId, _ := range subIdMap {
			if syntacticDistance[refId][subId] < smallestSyntacticDistance ||
				semanticMiniLM[refId][subId] > highestSemanticMiniScore {

				smallestSyntacticDistance = syntacticDistance[refId][subId]
				highestSemanticMiniScore = semanticMiniLM[refId][subId]
				highestScoreId = subId
			}
		}

		// fix this ID.
		refMMU := NewGraph()
		refMMU.Vertices[refId] = cmd.ReferenceSolution.Vertices[refId]
		subMMU := NewGraph()
		subMMU.Vertices[highestScoreId] = cmd.Submission.Vertices[highestScoreId]

		minimalUnits[&refMMU] = &subMMU
		fixedIds[refId] = highestScoreId

		// remove the ID from all other maps
		for _, subIdMap := range potentialMatches {
			delete(subIdMap, highestScoreId) // delete ID if the ID is present in any of the potential match maps
		}
	}

	// grow the equivalence graphs step-by-step
	progress := true
	for {
		if !progress {
			break
		}
		progress = false

		// make maximal combinations of fixed IDs (connect known minimal units together)
		for refMU, subMU := range minimalUnits {
			if refMU.Empty() || subMU.Empty() {
				delete(minimalUnits, refMU)
				continue
			}

			for refVId, _ := range refMU.Vertices {
				// get connected edges
				connectedEdges := cmd.ReferenceSolution.ConnectedEdges(refVId)

				// filter on which vertex IDs are fixed in the submission as well
				ind := 0
				connectedEdgesSub := []EdgeIdentifier{}
				for {
					if ind >= len(connectedEdges) {
						break
					}

					eId := connectedEdges[ind]
					// if the edge already appears in the MU, then it's not new, skip:
					if _, ok := refMU.Edges[eId]; ok {
						connectedEdges = slices.Delete(connectedEdges, ind, ind+1) // <- this is why we keep an index
						continue
					}

					e := cmd.ReferenceSolution.Edges[eId]
					otherV := e.FromId

					if e.FromId == nil || e.ToId == nil { // edge-edge connections are skipped
						connectedEdges = slices.Delete(connectedEdges, ind, ind+1) // <- this is why we keep an index
						continue
					}
					if *e.FromId == refVId {
						otherV = e.ToId
					}

					otherSubVId, ok := fixedIds[*otherV]
					connectedsubedges := cmd.Submission.ConnectedEdges(otherSubVId)
					if !ok || len(connectedsubedges) == 0 {
						connectedEdges = slices.Delete(connectedEdges, ind, ind+1) // <- this is why we keep an index
						continue
					}

					connectedEdgesSub = append(connectedEdgesSub, connectedsubedges...)
					ind++
				}

				if len(connectedEdges) != len(connectedEdgesSub) {
					c.LogErr("Trying to grow Minimal Units unequally!")
					return nil
				}
				if len(connectedEdges) == 0 {
					continue
				}
				progress = true

				// now grow the graphs by these edges
				refMU.MergeConnectedEdges(connectedEdges, cmd.ReferenceSolution)
				subMU.MergeConnectedEdges(connectedEdgesSub, cmd.Submission)
			}
			// ..and merge other graphs:
			/*
			 * For each other subgraph, if it conatins a vertex that is in this graph, then merge that graph into this one, including all edges etc.
			 */
			for g, s := range minimalUnits {
				if g != refMU {
					for vId, _ := range refMU.Vertices {
						if _, ok := g.Vertices[vId]; !ok {
							continue
						}

						// merge edge connections and their vertices into this graph:
						edges := g.ConnectedEdges(vId)
						refMU.MergeConnectedEdges(edges, g)
						g.DeleteEdgesAndTheirVertices(edges)
						delete(g.Vertices, vId) // if no connected edges, then we only have one vertex to delete
					}
				}
				if s != subMU {
					for vId, _ := range subMU.Vertices {
						if _, ok := s.Vertices[vId]; !ok {
							continue
						}

						edges := s.ConnectedEdges(vId)
						subMU.MergeConnectedEdges(edges, s)
						s.DeleteEdgesAndTheirVertices(edges)
						delete(s.Vertices, vId) // if no connected edges, then we only have one vertex to delete
					}
				}
			}
		}

		// try to explore edges of minimal units and discover more equivalences

	}

	return nil
}

func vertexToStr(v *InternalVertex) string {
	r := new(strings.Builder)
	r.WriteString(v.Title)
	r.WriteRune(' ')
	for n, v := range v.Values {
		r.WriteString(n)
		r.WriteRune(' ')
		valProps := []string{v.Value, v.Properties.Type, string(v.Properties.Visibility)}
		r.WriteString(strings.Join(valProps, ","))
		r.WriteRune(' ')
	}
	return r.String()
}

func getPotentialMatchingVertices(c *context.Ctx, cmd GradeCmd) (
	potentialMatches map[VertexIdentifier]map[VertexIdentifier]bool,
	syntacticDistance map[VertexIdentifier]map[VertexIdentifier]int,
	//semanticWordNet map[VertexIdentifier]map[VertexIdentifier]float64,
	semanticMiniLM map[VertexIdentifier]map[VertexIdentifier]float64,
) {
	potentialMatches = map[VertexIdentifier]map[VertexIdentifier]bool{}

	syntacticDistance = map[VertexIdentifier]map[VertexIdentifier]int{}
	// semanticWordNet = map[VertexIdentifier]map[VertexIdentifier]float64{}
	semanticMiniLM = map[VertexIdentifier]map[VertexIdentifier]float64{}

	// matching algorithm inspired by Smith et al. 2013 (forall vertex pairs, if syntactic match or semantic match, add to potential matches)

	for refId, refV := range cmd.ReferenceSolution.Vertices {
		syntacticDistance[refId] = map[VertexIdentifier]int{}
		// semanticWordNet[refId] = map[VertexIdentifier]float64{}
		semanticMiniLM[refId] = map[VertexIdentifier]float64{}

		for subId, subV := range cmd.Submission.Vertices {
			syntacticMatchCmd := MatchStringCmd{Ref: refV.Title, Act: subV.Title}
			syntDist := c.Queries.SyntacticMatch.Get("Syntactic match", syntacticMatchCmd)
			if syntDist <= MAX_ALLOWED_LEVENSHTEIN_DISTANCE {
				syntacticDistance[refId][subId] = syntDist
				if _, ok := potentialMatches[refId]; !ok {
					potentialMatches[refId] = map[VertexIdentifier]bool{}
				}
				potentialMatches[refId][subId] = true
				// continue
			}

			semanticMatchCmd := MatchStringCmd{Ref: vertexToStr(refV), Act: vertexToStr(subV)}

			resMini := c.Queries.SemanticMatchSentenceTransformer.Get("Semantic Match - MiniLM", semanticMatchCmd)
			if resMini.Err != nil {
				c.LogErr("Failed calculating similarity: %s", resMini.Err.Error())
				return nil, nil /*nil,*/, nil
			}

			//resWordNet := c.Queries.SemanticMatchWordnet.Get("Semantic Match - Wordnet", semanticMatchCmd)
			//if resWordNet.Err != nil {
			//	c.LogErr("Error while calculating semantic simlarity: %s", resWordNet.Err.Error())
			//	return nil
			//}

			//semanticWordNet[refId][subId] = resWordNet.Score
			if resMini.Score >= COSINE_SIMILARITY_THRESHOLD { //|| resWordNet.Score >= WORDNET_SIMILARITY_THRESHOLD {
				if _, ok := potentialMatches[refId]; !ok {
					potentialMatches[refId] = map[VertexIdentifier]bool{}
				}
				semanticMiniLM[refId][subId] = resMini.Score
				potentialMatches[refId][subId] = true
			}
		}
	}
	return potentialMatches, syntacticDistance /*semanticWordNet, */, semanticMiniLM
}
