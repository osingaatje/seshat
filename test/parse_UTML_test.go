package test

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/rogpeppe/go-internal/diff" // for diffing strings
)

func TestUTMLParseALLDATASETS(t *testing.T) {
	VerifyParsingForFiles(t, "") // uses 'double-star' GLOB syntax, supported by my matching function. Not in default Go!
}

// Tests whether the parsed and JSONified file and the formatted input file is literally, to the character, the exact same.
func TestUTMLParseSimple(t *testing.T) {
	VerifyParsingForFiles(t, "./examples/correct")
}

/********* CODE FOR MATCHING JSON DIFFS FOR PARSE RESULTS - IGNORING CERTAIN PROPERTIES *********/
// matches + or - "property": ...null/[]/"value" in JSON a.k.a. the diffed properties
var matchJSONDiffProperty = regexp.MustCompile(`([+-])\s+\"(.+)\":\s*(.+)\s*\n`) //`^\s*[+-]\s*\"(.+)\":\s*(.+)$`)
var ignoredProps []string = []string{
	"attributes", "methods",
}
var ignoredPropValues []string = []string{
	"null", "[]",
}
var matchProperty = func(match []string) (add bool, name string, value string) {
	// match: [entire string, attr, value (+ comma sometimes)]
	if len(match) != 4 {
		panic("REGEX MATCH INCORRECT WITH VALIDATING fUNCITON")
	}
	match[3] = strings.TrimSuffix(match[3], ",")

	// we ignore 'null'/'[]' diffs for 'attributes' or 'methods'
	return match[0] == "+", match[2], match[3]
}

/************************************************************************************************/

func VerifyParsingForFiles(t *testing.T, dir string, globs ...string) {
	c := driver.NewContext()

	var filepaths []string
	var err error
	if dir == "" && len(globs) == 0 {
		filepaths, err = helper.AllDatasetFiles()
	} else if len(globs) == 0 {
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

		// expectedJson, err := helper.RemarshalJSON(fileContents)
		// if err != nil {
		// 	t.Errorf("Could not remarshal JSON for file '%s' (to ensure sorted keys etc. Err=%s", path, err.Error())
		// 	return
		// }
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
		// resJSONBytes, err = helper.RemarshalJSON(resJSONBytes)
		// if err != nil {
		// 	t.Error("Could not REmarshal UTML parse result")
		// 	return
		// }
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

			differentProps := []string{}
			for _, m := range propertyMatches {
				_, prop, val := matchProperty(m)

				if slices.Contains(ignoredProps, prop) && slices.Contains(ignoredPropValues, val) {
					continue
				}
				// NON-IGNORED VALUES - skip if the value appears an even number of times (because of reordering, for ex.: "+attribute: val ..... -attribute: val")
				if i := slices.Index(differentProps, prop); i >= 0 {
					differentProps = slices.Delete(differentProps, i, i+1)
					continue
				}
				differentProps = append(differentProps, prop)
			}
			if len(differentProps) == 0 {
				continue
			}

			assert.Fail(t, "Different fields ["+strings.Join(differentProps, ",")+"] found in parsing: "+diffStr)
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
