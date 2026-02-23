package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	. "github.com/osingaatje/seshat/helper"
)

type MultTestCaseRes struct {
	Expected *Multiplicity
	ErrMsg   string
}

var testCases map[string]MultTestCaseRes = map[string]MultTestCaseRes{
	"":             {Expected: nil, ErrMsg: ""},
	"0":            {Expected: &Multiplicity{Start: 0, HasEndMult: false}, ErrMsg: ""},
	"0..1":         {Expected: &Multiplicity{Start: 0, HasEndMult: true, End: 1}, ErrMsg: ""},
	"1..15":        {Expected: &Multiplicity{Start: 1, HasEndMult: true, End: 15}, ErrMsg: ""},
	"*":            {Expected: &Multiplicity{Start: -1, HasEndMult: false}, ErrMsg: "'Many' should be translated to '-1'"},
	"*..*":         {Expected: &Multiplicity{Start: -1, HasEndMult: false}, ErrMsg: "Many to many should be normalised!"},
	"151230123..1": {Expected: &Multiplicity{Start: 1, HasEndMult: true, End: 151230123}, ErrMsg: "Inverse ranges should be normalised!"},
}

func TestParseMultiplicity(t *testing.T) {
	for testCase, exp := range testCases {
		actual, ok := GetMultiplicity(testCase)
		if exp.Expected == nil {
			assert.Equal(t, false, ok, "Nil result should have a 'false' ok value")
			return
		}

		assert.Equal(t, *exp.Expected, *actual, "Incorrectly parsed multiplicity '%s' into %+v (expected %+v) (errMsg: %s)", testCase, actual, exp.Expected, exp.ErrMsg)
	}
}
