package grade

import (
	"fmt"

	"github.com/fluhus/gostuff/nlp"
	"github.com/fluhus/gostuff/nlp/wordnet"
	wn "github.com/osingaatje/seshat/helper/wordnet"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
)

var net *wordnet.WordNet
var netErr error

func getMeanings(token string) map[wn.WordNetWordType][]*wordnet.Synset {
	meanings := net.Search(token)
	allMeanings := map[wn.WordNetWordType][]*wordnet.Synset{}
	for _, m := range wn.WordNetWordTypes {
		allMeanings[m] = meanings[string(m)]
	}
	return allMeanings
}

// Score: [0,1+]: 0 (totally not similar) - 1 (perfectly similar) (with a bit of leniency above 1 ("I'm 120% sure!!" type behaviour))
func semanticMatch(c *context.Ctx, s command.MatchStringCmd) command.MatchStringRes {
	// init variables etc.
	net, netErr = wn.GetWordNet()
	if netErr != nil {
		return command.MatchStringRes{
			Score: -1,
			Err:   fmt.Errorf("Could not get WordNet (for semantic matching), err=%s", netErr.Error()),
		}
	}

	var totalScore float64 = 0

	toks1 := nlp.Tokenize(s.Ref, false)
	toks2 := nlp.Tokenize(s.Act, false)

	// might look intimidating, but in real life, we mostly get single words or combinations of a couple words at most, with 4 word types, often containing 1 or 2 meanings.
	// so this might look like a n^5 function, and it is, but it's fine.
	for _, tok := range toks1 {
		tokMeanings := getMeanings(tok)
		for _, otherTok := range toks2 {
			otherTokMeanings := getMeanings(otherTok)

			for _, wordType := range wn.WordNetWordTypes {
				for _, meaning := range tokMeanings[wordType] {
					for _, otherMeaning := range otherTokMeanings[wordType] {
						totalScore += wn.Similarity(meaning, otherMeaning)
					}
				}
			}
		}
	}

	totalScore /= float64(len(toks1) + len(toks2)) // roughly normalise it to 0..1

	return command.MatchStringRes{
		Score: totalScore,
		Err:   nil,
	}
}
