package test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/rogpeppe/go-internal/diff" // for diffing strings
)

func TestUTMLParseALLDATASETS(t *testing.T) {
	VerifyParsingForFiles(t, "../DATASETS", "**/*.json") // uses 'double-star' GLOB syntax, supported by my matching function. Not in default Go!
}

// Tests whether the parsed and JSONified file and the formatted input file is literally, to the character, the exact same.
func TestUTMLParseSimple(t *testing.T) {
	VerifyParsingForFiles(t, "./examples/correct")
}

/********* CODE FOR MATCHING JSON DIFFS FOR PARSE RESULTS - IGNORING CERTAIN PROPERTIES *********/
// matches + or - "property": ...null/[]/"value" in JSON a.k.a. the diffed properties
var matchJSONDiffProperty = regexp.MustCompile(`[+-]\s+\"(.+)\":\s*(.+)\s*\n`) //`^\s*[+-]\s*\"(.+)\":\s*(.+)$`)
var matchIgnoredProperty = func(match []string) bool {
	// match: [entire string, attr, value (+ comma sometimes)]
	if len(match) != 3 {
		panic("REGEX MATCH INCORRECT WITH VALIDATING fUNCITON")
	}
	match[2] = strings.TrimSuffix(match[2], ",")

	// we ignore 'null'/'[]' diffs for 'attributes' or 'methods'
	return (match[1] == "attributes" || match[1] == "methods") && (match[2] == "null" || match[2] == "[]")
}

/************************************************************************************************/

func VerifyParsingForFiles(t *testing.T, dir string, globs ...string) {
	c := driver.NewContext()

	var filepaths []string
	var err error
	if len(globs) == 0 {
		filepaths, err = helper.AllUTMLFiles(dir)
	} else {
		filepaths, err = helper.AllFiles(dir, globs...)
	}
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("Matched %d files for verifying parsing. (Base dir='%s')", len(filepaths), dir)

	for _, path := range filepaths {

		fileContents, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Could not read file: Err=%s", err.Error())
			return
		}

		expectedJson, err := helper.IndentJSON(fileContents)
		if err != nil {
			t.Errorf("Could not indent JSON for file '%s', err=%s", path, err.Error())
			return
		}

		// run query
		res := c.Queries.ParseUTML.Get("Parse UTML", path)
		if res == nil {
			t.Error("Could not parse UTML!")
			return
		}

		// convert the parse result into an 'any' the same way to avoid type comparison
		resJSONBytes, err := helper.MarshalJSON(res)
		if err != nil {
			t.Error("Could not marshal UTML parse result")
			return
		}
		resJSONBytes, err = helper.IndentJSON(resJSONBytes)
		if err != nil {
			t.Error("Could not indent UTML parse result")
			return
		}

		diffs := diff.Diff("Expected", expectedJson, "Actual", resJSONBytes)
		if diffs != nil {
			// hacky shit: if we detect that "attributes" or "methods" have been added with 'null' value, ignore those as well as  '}' -> '},'.
			// I could not find a more nice method than this (trust me, I tried `co-cmp` and everything, but it all turned into a hacky mess anyway.)
			diffStr := string(diffs)
			propertyMatches := matchJSONDiffProperty.FindAllStringSubmatch(diffStr, 1000)

			hasRelevantDiff := false
			for _, m := range propertyMatches {
				if !matchIgnoredProperty(m) {
					hasRelevantDiff = true
					break
				}
			}
			if !hasRelevantDiff {
				continue
			}

			assert.Fail(t, "Difference found in parsing:", diffStr)
		}
	}
}

func TestUTMLBroken(t *testing.T) {

	var filePaths []string = helper.AllUTMLFilesUNSAFE("./examples/broken-utml")

	for _, path := range filePaths {
		// init fresh context without existing logs
		c := driver.NewContext()

		// run query
		res := c.Queries.ParseUTML.Get("Parse UTML", path)

		assert.Nil(t, res) // query should fail
		if res != nil {
			r, err := helper.MarshalJSONWithIndent(res)
			if err != nil {
				return
			}
			t.Logf("Produced graph: \n%s", r)
		}

		strings.Contains(c.Logger.GetLogString(), "Could not marshal file")
	}
}
