package grade

import (
	"math"
	"strings"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	internalrep "github.com/osingaatje/seshat/types/graph/internal-rep"
	"github.com/osingaatje/seshat/types/graph/shared"
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

	potentialMatches := map[shared.VertexIdentifier]map[shared.VertexIdentifier]bool{}

	syntacticDistance := map[shared.VertexIdentifier]map[shared.VertexIdentifier]int{}
	semanticWordNet := map[shared.VertexIdentifier]map[shared.VertexIdentifier]float64{}
	semanticMiniLM := map[shared.VertexIdentifier]map[shared.VertexIdentifier]float64{}

	// matching algorithm inspired by Smith et al. 2013 (forall vertex pairs, if syntactic match or semantic match, add to potential matches)

	for refId, refV := range cmd.ReferenceSolution.Vertices {
		syntacticDistance[refId] = map[shared.VertexIdentifier]int{}
		semanticWordNet[refId] = map[shared.VertexIdentifier]float64{}
		semanticMiniLM[refId] = map[shared.VertexIdentifier]float64{}

		for subId, subV := range cmd.Submission.Vertices {
			syntacticMatchCmd := MatchStringCmd{Ref: refV.Title, Act: subV.Title}
			syntacticDistance[refId][subId] = c.Queries.SyntacticMatch.Get("Syntactic match", syntacticMatchCmd)
			if syntacticDistance[refId][subId] <= MAX_ALLOWED_LEVENSHTEIN_DISTANCE {
				if _, ok := potentialMatches[refId]; !ok {
					potentialMatches[refId] = map[shared.VertexIdentifier]bool{}
				}
				potentialMatches[refId][subId] = true
				// continue
			}

			semanticMatchCmd := MatchStringCmd{Ref: vertexToStr(refV), Act: vertexToStr(subV)}

			resMini := c.Queries.SemanticMatchSentenceTransformer.Get("Semantic Match - MiniLM", semanticMatchCmd)
			if resMini.Err != nil {
				c.LogErr("Failed calculating similarity: %s", resMini.Err.Error())
				return nil
			}

			//resWordNet := c.Queries.SemanticMatchWordnet.Get("Semantic Match - Wordnet", semanticMatchCmd)
			//if resWordNet.Err != nil {
			//	c.LogErr("Error while calculating semantic simlarity: %s", resWordNet.Err.Error())
			//	return nil
			//}

			semanticMiniLM[refId][subId] = resMini.Score
			//semanticWordNet[refId][subId] = resWordNet.Score
			if resMini.Score >= COSINE_SIMILARITY_THRESHOLD { //|| resWordNet.Score >= WORDNET_SIMILARITY_THRESHOLD {
				potentialMatches[refId][subId] = true
			}
		}
	}

	fixedIds := map[shared.VertexIdentifier]shared.VertexIdentifier{}

	for refId, subIdMap := range potentialMatches {
		var smallestSyntacticDistance int = math.MaxInt
		var highestSemanticMiniScore float64 = math.MaxFloat64
		var highestScoreId shared.VertexIdentifier = 0

		for subId, _ := range subIdMap {
			if syntacticDistance[refId][subId] < smallestSyntacticDistance ||
				semanticMiniLM[refId][subId] > highestSemanticMiniScore {

				smallestSyntacticDistance = syntacticDistance[refId][subId]
				highestSemanticMiniScore = semanticMiniLM[refId][subId]
				highestScoreId = subId
			}
		}

		// fix this ID.
		fixedIds[refId] = highestScoreId
	}

	//if len(fixedIds) == len(cmd.ReferenceSolution.Vertices) && len(cmd.ReferenceSolution.Vertices) == len(cmd.Submission.Vertices) {
	//	// we fixed everything I guess?
	//
	//}
	return nil
}

func vertexToStr(v *internalrep.InternalVertex) string {
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
