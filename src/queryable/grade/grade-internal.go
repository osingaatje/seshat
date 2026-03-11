package grade

import (
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/fluhus/gostuff/nlp"
	"github.com/fluhus/gostuff/nlp/wordnet"
	wn "github.com/osingaatje/seshat/helper/wordnet"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
)

var net *wordnet.WordNet
var netErr error

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

	// init variables etc.
	net, netErr = wn.GetWordNet()
	if netErr != nil {
		c.LogErr("Could not get WordNet (for semantic matching), err=%s", netErr.Error())
		return nil
	}

	// certainties := map[graphVertexIdentifier]
	//
	//	for id, v := range cmd.ReferenceSolution.Vertices {
	//
	//	}

	return nil
}

func semanticMatch(str1 string, str2 string) {
	if net == nil || netErr != nil {
		panic("WordNet not correctly initialised!")
	}

	var totalScore float64 = 0

	toks1 := nlp.Tokenize(str1, false)
	toks2 := nlp.Tokenize(str2, false)

	for _, tok := range toks1 {
		tokMeanings := getMeanings(tok)
		for _, otherTok := range toks2 {
			otherTokMeanings := getMeanings(otherTok)

			for _, meaning := range tokMeanings {
				for _, otherMeaning := range otherTokMeanings {
					totalScore += wn.Similarity(meaning, otherMeaning)
				}
			}
		}
	}
}
func getMeanings(token string) []*wordnet.Synset {
	meanings := net.Search(token)
	allMeanings := []*wordnet.Synset{}
	for _, m := range wn.WordNetWordTypes {
		allMeanings = append(allMeanings, meanings[string(m)]...)
	}
	return allMeanings
}

// Computes Levenshtein distance between normalised strings
func syntacticMatch(str1 string, str2 string) int {
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
