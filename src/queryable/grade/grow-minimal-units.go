package grade

import (
	"math"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/graph/internal-rep"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

const CERTAINTY_THRESHOLD float64 = 0.7 // applicable to the combination of semantic and syntactic scores ('certaintyMap')

/*
 * Idea (inspired by Smith (2004), Thomas (2009), Vachharajani (2014), Bian (2020))
 * (steps inspired by Smith 2004, skipping segmentation and assimilation because the internal graph repr. already does that)
 *
 * 1. Identify / match: using (Minimal) Meaningful Units, with algorithms such as levenshtein distance (syntactic) and WordNet similarity score (HSO/WUP/LIN) (semantic)
 * 2. Aggregate: 		continually increase Meaningful Unit until it can no longer be aggregated) using structural graph matching
 * 3. Interpret:        produce grades based on present/missing/extra elements and parts of those elements
 */
func getAlternativeSolutions(
	c *context.Ctx, cmd GradeCmd) []map[VertexIdentifier]VertexIdentifier {
	if cmd.ReferenceSolution == nil || cmd.Submission == nil || cmd.Rubric == nil {
		c.LogErr("Referencesolution, submission, or rubric was not present when grading diagram!")
		return nil
	}

	potentialMatches, syntacticDistance, semanticMiniLM := getPotentialMatchingVertices(c, cmd)
	certaintyMap := combineSyntacticSemanticScore(cmd.ReferenceSolution, cmd.Submission, syntacticDistance, semanticMiniLM)

	/*
	 * This map contains the minimal meaningful units, mapped from REFERENCE -> SUBMISSION
	 */
	minimalUnits := map[*InternalGraph]*InternalGraph{}
	// per graph->graph matching: a map of vertex->vertex
	fixedVertices := map[*InternalGraph]map[VertexIdentifier]VertexIdentifier{}

	progresses := map[*InternalGraph]bool{} // if the graph was further expanded this refinement step

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

		// skip unsimilar vertices
		if smallestSyntacticDistance > SYNTACTIC_DISTANCE_THRESHOLD && highestSemanticMiniScore < COSINE_SIMILARITY_THRESHOLD {
			continue
		}

		// fix this ID.
		refMMU := NewGraph()
		refMMU.Vertices[refId] = cmd.ReferenceSolution.Vertices[refId]
		subMMU := NewGraph()
		subMMU.Vertices[highestScoreId] = cmd.Submission.Vertices[highestScoreId]

		minimalUnits[&refMMU] = &subMMU

		fixedVertices[&refMMU] = map[VertexIdentifier]VertexIdentifier{}
		fixedVertices[&refMMU][refId] = highestScoreId

		progresses[&refMMU] = true

		// remove the ID from all other maps
		for _, subIdMap := range potentialMatches {
			delete(subIdMap, highestScoreId) // delete ID if the ID is present in any of the potential match maps
		}
	}

	// grow the equivalence graphs step-by-step
	for {
		// if there's not any progress, then stop
		if !helper.AnyMap(progresses) {
			break
		}

		for refMU, subMU := range minimalUnits {
			if !progresses[refMU] { // skip already finished expansions
				continue
			}
			progress := growUnits(cmd.ReferenceSolution, cmd.Submission, refMU, subMU, certaintyMap, fixedVertices[refMU])
			progresses[refMU] = progress
		}
	}

	return helper.ValuesMap(fixedVertices)
}

/*
 * Idea:
 * 1. Explore connections in the reference and submission MU for all vertices
 * 2. Filter out edges/vertices that are already in the MUs
 *     -> Choose the most fitting vertices to connect based on the certainty map (given that the match is above the threshold)
 * 3. Check if we can add any new vertices / edges (if so, progress!)
 *
 */
func growUnits(
	refGraph *InternalGraph,
	subGraph *InternalGraph,
	refMU *InternalGraph,
	subMU *InternalGraph,
	certaintyMap map[VertexIdentifier]map[VertexIdentifier]float64,
	fixedIds map[VertexIdentifier]VertexIdentifier) /*progress*/ bool {

	newEdgeRefs, newVertRefs := getNewEdges(refGraph, refMU)
	newEdgeSubs, newVertSubs := getNewEdges(subGraph, subMU)
	if len(newEdgeRefs) == 0 && len(newEdgeSubs) == 0 {
		return false
	}

	newFixedIds := map[VertexIdentifier]VertexIdentifier{}

	// fix new vertex to vertex correlations:
	// TODO: potentially sort this by the most certain match across all vertices (more complicated though)
	for vId, _ := range newVertRefs {
		if _, ok := fixedIds[vId]; ok {
			continue
		}

		var bestMatchSubId VertexIdentifier
		var bestMatchScore float64 = -1

		// if we have assigned all reference/submission
		if len(newVertSubs)-len(newFixedIds) <= 0 || len(newVertRefs)-len(newFixedIds) <= 0 {
			break
		}

		for subVId, _ := range newVertSubs {
			// don't reassign already given out vertex IDs
			if slices.Contains(helper.ValuesMap(fixedIds), subVId) {
				continue
			}

			if certaintyMap[vId][subVId] > bestMatchScore {
				bestMatchScore = certaintyMap[vId][subVId]
				bestMatchSubId = subVId
			}
		}
		if bestMatchScore <= 0 {
			continue
		}

		newFixedIds[vId] = bestMatchSubId
	}

	// for each fixed ID, add the new edges and vertices, ONLY IF THEY APPEAR IN THE FIXED IDS MAP
	for refId, subId := range newFixedIds {
		refMU.Vertices[refId] = refGraph.Vertices[refId]
		subMU.Vertices[subId] = subGraph.Vertices[subId]
		fixedIds[refId] = subId
	}

	// only add new edges that have a to/from ID in the fixedIds
	for refEId, refE := range newEdgeRefs {
		if refE.FromId == nil || refE.ToId == nil {
			continue
		}
		if _, ok := fixedIds[*refE.FromId]; !ok {
			continue
		}
		if _, ok := fixedIds[*refE.ToId]; !ok {
			continue
		}
		// now add the edge
		refMU.Edges[refEId] = refE
	}

	for subEId, subE := range newEdgeSubs {
		if subE.FromId == nil || subE.ToId == nil {
			continue
		}
		if _, ok := fixedIds[*subE.FromId]; !ok {
			continue
		}
		if _, ok := fixedIds[*subE.ToId]; !ok {
			continue
		}
		// now add the edge
		subMU.Edges[subEId] = subE
	}

	return true
}
func getNewEdges(originalGraph *InternalGraph, mu *InternalGraph) (map[EdgeIdentifier]*InternalEdge, map[VertexIdentifier]*InternalVertex) {
	newEdges := map[EdgeIdentifier]*InternalEdge{}
	newVerts := map[VertexIdentifier]*InternalVertex{}
	for vId, _ := range mu.Vertices {
		newEIds := originalGraph.ConnectedEdges(vId)
		for _, eId := range newEIds {
			if _, ok := mu.Edges[eId]; ok { // add only edges that don't appear already in the graph
				continue
			}

			edge := originalGraph.Edges[eId]
			newEdges[eId] = edge
			if edge.FromId != nil && mu.Vertices[*edge.FromId] == nil { // if we havea fromId and it's not in the MU
				newVerts[*edge.FromId] = originalGraph.Vertices[*edge.FromId]
			}
			if edge.ToId != nil && mu.Vertices[*edge.ToId] == nil { // if we have a ToId and it's not in the MU
				newVerts[*edge.ToId] = originalGraph.Vertices[*edge.ToId]
			}
		}
	}

	return newEdges, newVerts
}

func combineSyntacticSemanticScore(
	refGraph *InternalGraph,
	subGraph *InternalGraph,
	syntDist map[VertexIdentifier]map[VertexIdentifier]int,
	semDist map[VertexIdentifier]map[VertexIdentifier]float64) map[VertexIdentifier]map[VertexIdentifier]float64 {

	res := map[VertexIdentifier]map[VertexIdentifier]float64{}

	for refVertex, syntMap := range syntDist {
		semMap, ok := semDist[refVertex]
		if !ok {
			panic("bug in the code, syntactic and semantic distance maps unequal")
		}
		res[refVertex] = map[VertexIdentifier]float64{}

		for subVertex, syntV := range syntMap {
			semV, ok := semMap[subVertex]
			if !ok {
				panic("bug in code, semantic/syntactic submission maps unequal")
			}

			refV, okR := refGraph.Vertices[refVertex]
			subV, okS := subGraph.Vertices[subVertex]
			if !okR || !okS {
				panic("bug in code, certainty for a nonexisting vertex pair")
			}

			res[refVertex][subVertex] = combineScores(refV, subV, syntV, semV)
		}
	}
	return res
}

func combineScores(refV *InternalVertex, subV *InternalVertex, syn int, sem float64) float64 {
	if syn == 0 { // exactly equal strings get special treatment :)
		return 1
	}
	//cursed but nice oneliner
	longestWord := math.Max(float64(len(vertexToStr(refV))), float64(len(vertexToStr(subV))))

	/*
			 *  (1 - syn/wordlength) + sem
			 *   ------------------------
		     *              2
	*/
	return ((1 - (float64(syn) / longestWord)) + sem) / 2
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
			syntacticDistance[refId][subId] = syntDist
			if syntDist <= SYNTACTIC_DISTANCE_THRESHOLD {
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

			semanticMiniLM[refId][subId] = resMini.Score
			//semanticWordNet[refId][subId] = resWordNet.Score
			if resMini.Score >= COSINE_SIMILARITY_THRESHOLD { //|| resWordNet.Score >= WORDNET_SIMILARITY_THRESHOLD {
				if _, ok := potentialMatches[refId]; !ok {
					potentialMatches[refId] = map[VertexIdentifier]bool{}
				}
				potentialMatches[refId][subId] = true
			}
		}
	}
	return potentialMatches, syntacticDistance /*semanticWordNet, */, semanticMiniLM
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
