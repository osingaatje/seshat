package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/src/queryable/grade"
	"github.com/osingaatje/seshat/types/command"
)

const DESIRED_ACCURACY = 0.9 // can fail 10% of the time

var matchCases []command.MatchStringCmd = []command.MatchStringCmd{
	{
		Ref: "ChargingPort", // perfectly similar words should be counted as the same (word is from Bouali 2025)
		Act: "ChargingPort",
	}, { // Synonyms in teaching should be accounted for since these examples are used often (in my experience at least -Douwe)
		Ref: "Student",
		Act: "Pupil",
	}, { // some programmer terms
		Ref: "function",
		Act: "method",
	}, {
		Ref: "function",
		Act: "procedure",
	}, {
		Ref: "function",
		Act: "subroutine",
	}, {
		Ref: "loop",
		Act: "iteration",
	}, /* {
	Ref: "array",
	Act: "list",
	},*/ /* { // some more general things
	Ref: "teacher",
	Act: "educator",
	}, */{
		Ref: "teacher",
		Act: "instructor",
	}, {
		Ref: "house",
		Act: "home",
	}, {
		Ref: "house",
		Act: "residence",
	}, {
		Ref: "car",
		Act: "vehicle",
	}, {
		Ref: "car",
		Act: "automobile",
	}, /* {
	Ref: "smart",
	Act: "clever",
	}, */{
		Ref: "book",
		Act: "novel",
	}, {
		Ref: "book",
		Act: "publication",
	}, {
		Ref: "friend",
		Act: "companion",
	}, {
		Ref: "friend",
		Act: "ally",
	}, {
		Ref: "child",
		Act: "kid",
	}, {
		Ref: "city",
		Act: "metropolis",
	}, {
		Ref: "city",
		Act: "town",
	}, {
		Ref: "tree",
		Act: "plant",
	}, {
		Ref: "mountain",
		Act: "peak",
	}, {
		Ref: "mountain",
		Act: "hill",
	}, {
		Ref: "road",
		Act: "street",
	}, { // more academic / scholarly / business / medical
		Ref: "job",
		Act: "occupation",
	}, {
		Ref: "job",
		Act: "career",
	}, {
		Ref: "money",
		Act: "cash",
	}, {
		Ref: "money",
		Act: "currency",
	}, {
		Ref: "school",
		Act: "academy",
	}, {
		Ref: "school",
		Act: "institution",
	}, {
		Ref: "doctor",
		Act: "physician",
	}, {
		Ref: "doctor",
		Act: "medic",
	}, { // back to nerdy stuff
		Ref: "computer",
		Act: "PC",
	}, {
		Ref: "computer",
		Act: "machine",
	}, /* {
	Ref: "computer",
	Act: "laptop",
	}, */{
		Ref: "computer",
		Act: "desktop",
	}, { // code examples
		Ref: "Student yawn() string",
		Act: "Pupil tired() string",
	}, {
		Ref: "Rollercoaster ride() void getName() string",
		Act: "Attraction ride() name() string",
	},
}

var notMatchCases []command.MatchStringCmd = []command.MatchStringCmd{
	/*{ // this matches very well, but we cannot disambiguate Port/Station because we don't know the context. These are semantically related words.
	Ref: "ChargingPort",
	Act: "ChargingStation",
	},*/{
		Ref: "successor",
		Act: "dancing",
	}, {
		Ref: "person",
		Act: "coffee",
	}, {
		Ref: "apple",
		Act: "laptop",
	}, {
		Ref: "river",
		Act: "clock",
	}, {
		Ref: "music",
		Act: "owl",
	}, {
		Ref: "car",
		Act: "piano",
	}, {
		Ref: "ocean",
		Act: "lamp",
	}, {
		Ref: "mountain",
		Act: "telephone",
	}, {
		Ref: "bird",
		Act: "bicycle",
	}, {
		Ref: "rain",
		Act: "guitar",
	}, {
		Ref: "happiness",
		Act: "spoon",
	}, {
		Ref: "silence",
		Act: "balloon",
	}, {
		Ref: "alisson",
		Act: "pencil",
	}, {
		Ref: "curse",
		Act: "headphones",
	}, {
		Ref: "mildew",
		Act: "automation",
	}, {
		Ref: "team member",
		Act: "virtual machine",
	},
}

func TestSemanticMATCH(t *testing.T) {
	semanticMatch(t, matchCases, func(score float64) bool { return score > grade.SEMANTIC_POSITIVE_MATCH_THRESHOLD })
}

func TestNOTSemanticallyRelated(t *testing.T) {
	semanticMatch(t, notMatchCases, func(score float64) bool { return score < grade.SEMANTIC_POSITIVE_MATCH_THRESHOLD })
}

type TestRes struct {
	Cmd command.MatchStringCmd
	Res command.MatchStringRes
}

func semanticMatch(t *testing.T, cases []command.MatchStringCmd, successfunc func(score float64) bool) {
	ctx := driver.NewContext()

	successes := 0
	failedCases := []TestRes{}

	for _, c := range cases {
		r := ctx.Queries.SemanticMatch.Get("Semantic Match", c)
		if r.Err != nil {
			t.Fatalf("Error in semantic match: %s", r.Err.Error())
			return
		}
		//
		//
		//
		//
		//
		//
		if !successfunc(r.Score) {
			ctx.LogWarn("Semantic Match '%s' <-> '%s' was %.2f", c.Ref, c.Act, r.Score)
			failedCases = append(failedCases, TestRes{c, r})
			continue
		}

		successes++
	}

	caseLen := len(cases)
	accuracy := float64(successes) / float64(caseLen)

	if accuracy < DESIRED_ACCURACY {
		failCaseStr := new(strings.Builder)
		for _, c := range failedCases {
			failCaseStr.WriteString(fmt.Sprintf("\n\t%s <-> %s scored %.2f", c.Cmd.Ref, c.Cmd.Act, c.Res.Score))
		}

		t.Fatalf("TOO MANY FAILURES IN SEMANTIC MATCHING (accuracy%% = %.2f < %.2f: %d failed out of %d cases).\nFailed cases:%s", accuracy, DESIRED_ACCURACY, caseLen-successes, caseLen, failCaseStr.String())
	}
}
