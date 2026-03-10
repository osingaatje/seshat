package word

import (
	"fmt"

	"github.com/fluhus/gostuff/nlp/wordnet"
)

type WordNetWordType string

const (
	WordNetAdjective WordNetWordType = "a"
	WordNetNoun      WordNetWordType = "n"
	WordNetAdverb    WordNetWordType = "r"
	WordNetVerb      WordNetWordType = "v"
)

var WordNetWordTypes []WordNetWordType = []WordNetWordType{
	WordNetAdjective,
	WordNetNoun,
	WordNetAdverb,
	WordNetVerb,
}

var wn *wordnet.WordNet = nil
var err error = nil

func GetWordNet() (*wordnet.WordNet, error) {
	if err != nil {
		return nil, err
	}
	if wn != nil {
		return wn, nil
	}

	wn, err = wordnet.Parse("./dict")
	if err != nil {
		return nil, fmt.Errorf("Error occurred while reading WordNet: %s", err.Error())
	}
	return wn, nil
}

func Similarity(s1 *wordnet.Synset, s2 *wordnet.Synset) float64 {
	return wn.PathSimilarity(s1, s2, false)
}
