package grade

import (
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
)

const SYNTACTIC_DISTANCE_THRESHOLD = 2 // no. of characters

// Computes Levenshtein distance between normalised strings
func syntacticDist(c *context.Ctx, cmd MatchStringCmd) int {
	s1 := removePunctuation(cmd.Ref)
	s2 := removePunctuation(cmd.Act)

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
