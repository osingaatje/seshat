package grade

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	lemma "github.com/aaaton/golem/v4"
	lemma_en "github.com/aaaton/golem/v4/dicts/en"
	"github.com/osingaatje/gostuff/nlp/wordnet"
	wn "github.com/osingaatje/seshat/helper/wordnet"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
)

var (
	net    *wordnet.WordNet
	netErr error
	lem    *lemma.Lemmatizer
	lemErr error

	spaces *regexp.Regexp = regexp.MustCompile(`\s+`)
)

func getMeanings(token string) map[wn.WordNetWordType][]*wordnet.Synset {
	meanings := net.Search(token)
	allMeanings := map[wn.WordNetWordType][]*wordnet.Synset{}
	for _, m := range wn.WordNetWordTypes {
		allMeanings[m] = meanings[string(m)]
	}
	return allMeanings
}

const SEMANTIC_POSITIVE_MATCH_THRESHOLD float64 = 0.5

// Score: [0,1+]: 0 (totally not similar) - 1 (perfectly similar) (with a bit of leniency above 1 ("I'm 120% sure!!" type behaviour))
func semanticMatch(c *context.Ctx, s command.MatchStringCmd) command.MatchStringRes {
	if strings.EqualFold(s.Ref, s.Act) { // if strings are case-insentively equal, return 1
		return command.MatchStringRes{
			Score: 1,
			Err:   nil,
		}
	}

	// init variables etc.
	if net == nil && netErr == nil {
		net, netErr = wn.GetWordNet()
	}
	if netErr != nil {
		return command.MatchStringRes{
			Score: -1,
			Err:   fmt.Errorf("Could not get WordNet (for semantic matching), err=%s", netErr.Error()),
		}
	}

	if lem == nil && lemErr == nil {
		lem, lemErr = lemma.New(lemma_en.New())
	}

	if lemErr != nil {
		return command.MatchStringRes{
			Score: -1,
			Err:   fmt.Errorf("Could not get Lemmatiser (for semantic matching), err=%s", lemErr.Error()),
		}
	}

	var totalScore float64 = 0

	refTokens := splitUpString(s.Ref)
	actTokens := splitUpString(s.Act)

	for i, r := range refTokens {
		refTokens[i] = lem.Lemma(r)
	}
	for i, a := range actTokens {
		actTokens[i] = lem.Lemma(a)
	}

	// might look intimidating, but in real life, we mostly get single words or combinations of a couple words at most, with 4 word types, often containing 1 or 2 meanings.
	// so this might look like a n^5 function, and it is, but it's fine.
	for _, tok := range refTokens {
		tokMeanings := getMeanings(tok)
		for _, otherTok := range actTokens {
			otherTokMeanings := getMeanings(otherTok)

			// compare strictly nouns to nouns, verbs to verbs (less expensive, more strict)
			//for _, wordType := range wn.WordNetWordTypes {
			//	for _, meaning := range tokMeanings[wordType] {
			//		for _, otherMeaning := range otherTokMeanings[wordType] {
			//			totalScore += wn.Similarity(meaning, otherMeaning)
			//		}
			//	}
			//}

			// get all meanings (nouns, verbs, you name it) and compare against all other meanings (more computationally expensive)
			allTokMeanings := []*wordnet.Synset{}
			allOtherTokMeanings := []*wordnet.Synset{}
			for _, wordType := range wn.WordNetWordTypes {
				allTokMeanings = append(allTokMeanings, tokMeanings[wordType]...)
				allOtherTokMeanings = append(allOtherTokMeanings, otherTokMeanings[wordType]...)
			}

			// check all meanings with each other
			for _, meaning := range allTokMeanings {
				for _, otherMeaning := range allOtherTokMeanings {
					totalScore += wn.Similarity(meaning, otherMeaning)
				}
			}
		}
	}

	totalScore /= float64(len(refTokens) + len(actTokens)) // roughly normalise it to 0..1

	return command.MatchStringRes{
		Score: totalScore,
		Err:   nil,
	}
}

/*
 * Split up string according to general programmer habits:
 * - camelCase, PascalCase, snake_case, kebab-case, etc.
 * - lower all characters (A -> a)
 *
 * For example:
 * - "TestObject" becomes ["test", "object"] and gets returned as "test object"
 * - "test_object" also becomes "test object"
 */
func splitUpString(s string) []string {
	res := new(strings.Builder)

	isLikelyCapsImportant := false
	if len(s) >= 2 && unicode.IsUpper(rune(s[0])) != unicode.IsUpper(rune(s[1])) || !strings.ContainsAny(s, "-_") {
		isLikelyCapsImportant = true
	}

	for _, c := range s {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			res.WriteRune(' ')
		}
		if isLikelyCapsImportant && unicode.IsUpper(c) && res.Len() > 0 {
			res.WriteRune(' ')
		}

		res.WriteRune(unicode.ToLower(c))
	}

	const SEP = "|!|"
	strRes := spaces.ReplaceAllString(res.String(), SEP)
	return strings.Split(strRes, SEP)
}
