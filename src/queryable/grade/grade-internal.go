package grade

import (
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
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

	syntacticDistance := map[shared.VertexIdentifier]map[shared.VertexIdentifier]int{}
	certainties := map[shared.VertexIdentifier]map[shared.VertexIdentifier]float64{}

	for refId, refV := range cmd.ReferenceSolution.Vertices {
		syntacticDistance[refId] = map[shared.VertexIdentifier]int{}
		certainties[refId] = map[shared.VertexIdentifier]float64{}

		for subId, subV := range cmd.Submission.Vertices {
			syntacticDistance[refId][subId] = syntacticDist(refV.Title, subV.Title)

			res := c.Queries.SemanticMatch.Get("Semantic Match", MatchStringCmd{Ref: refV.Title, Act: subV.Title})

			if res.Err != nil {
				c.LogErr("Error while calculating semantic simlarity: %s", res.Err.Error())
				return nil
			}
			certainties[refId][subId] = res.Score
		}
	}

	c.LogErr("TODO GRADE FURTHER")
	return nil
}

// Computes Levenshtein distance between normalised strings
func syntacticDist(str1 string, str2 string) int {
	s1 := removePunctuation(str1)
	s2 := removePunctuation(str2)

	return levenshtein.ComputeDistance(s1, s2)
}

// inspired by https://stackoverflow.com/questions/32081808/strip-all-whitespace-from-a-string
func removePunctuation(str string) string {
	b := new(strings.Builder)
	b.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
