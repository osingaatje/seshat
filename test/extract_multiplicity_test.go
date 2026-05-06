package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	. "github.com/osingaatje/seshat/helper/multiplicity"
)

type MultTestCaseRes struct {
	HasMult  bool // default = false
	Expected Multiplicity
	ErrMsg   string
}

var EmptyMult = Multiplicity{}

var testCases map[string]MultTestCaseRes = map[string]MultTestCaseRes{
	"":      {ErrMsg: ""},
	"0":     {HasMult: true, Expected: Multiplicity{Start: 0, HasEndMult: false}, ErrMsg: ""},
	"0..1":  {HasMult: true, Expected: Multiplicity{Start: 0, HasEndMult: true, End: 1}, ErrMsg: ""},
	"1..15": {HasMult: true, Expected: Multiplicity{Start: 1, HasEndMult: true, End: 15}, ErrMsg: ""},
	"*":/* treated as 0..* */ {HasMult: true, Expected: Multiplicity{Start: 0, HasEndMult: true, End: -1}, ErrMsg: "'Many' should be translated to '0..*'"},
	"*..*":/* also treated as 0..* */ {HasMult: true, Expected: Multiplicity{Start: 0, HasEndMult: true, End: -1}, ErrMsg: "Many to many should be normalised!"},
	"151230123..1":                   {HasMult: true, Expected: Multiplicity{Start: 1, HasEndMult: true, End: 151230123}, ErrMsg: "Inverse ranges should be normalised!"},
	"somethingthatcannotbeparsed..1": {ErrMsg: "Raw text should not be parseable!"},
	"-=-=-=- =-=-":                   {ErrMsg: "Raw text should not be parseable!"},
}

func TestParseMultiplicity(t *testing.T) {
	for testCase, exp := range testCases {
		actual, ok := GetMultiplicity(testCase)
		if !exp.HasMult {
			assert.Equal(t, false, ok, "Nil result should have a 'false' ok value")
			return
		}

		assert.Equal(t, exp.Expected, actual, "Incorrectly parsed multiplicity '%s' into %+v (expected %+v) (errMsg: %s)", testCase, actual, exp.Expected, exp.ErrMsg)
	}
}
