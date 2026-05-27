package grade

import (
	sentencetransformer "github.com/osingaatje/seshat/helper/sentence-transformer"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
)

func semanticSimilarityMiniLM(c *context.Ctx, in MatchStringCmd) MatchStringRes {
	r, err := sentencetransformer.CompareSentences(c, []string{in.Ref, in.Act})
	if err != nil {
		return MatchStringRes{
			Score: -1,
			Err:   err,
		}
	}
	if len(r) != 2 {
		return MatchStringRes{
			Score: -1,
			Err:   c.LogErrAndReturn("Did not expect more or less than two result vectors from comparing two strings!"),
		}
	}

	sim := sentencetransformer.CosineSimilarity(r[0], r[1])
	//	if sim < -1*math.Pi || sim > math.Pi { <-- Debug stuff
	//		c.LogWarn("Cosine should not go outside [-pi,pi]! Cosine sim: %.2f", sim)
	//	}
	return MatchStringRes{
		Score: sim,
		Err:   nil,
	}
}
