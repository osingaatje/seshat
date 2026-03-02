package test

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/test/helpers"
	"github.com/osingaatje/seshat/types/command"
)

// Tests whether the parsed and JSONified file and the formatted input file is literally, to the character, the exact same.
func TestUTMLParseSimple(t *testing.T) {
	c := driver.NewContext()

	var FilePaths []string = helpers.AllUTMLFiles("./examples")

	for _, path := range FilePaths {

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
		res := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: path})
		actualJson, err := helper.IndentJSON([]byte(res.String()))
		if err != nil {
			t.Errorf("Could not indent JSON for parsed UTML, err=%s", err.Error())
			return
		}

		// Write to temp file - useful to diff the JSON in an actual editor.
		//	err = os.WriteFile("test.utml", []byte(actualJson), os.ModePerm)
		//	if err != nil {
		//		t.FailNow()
		//	}

		assert.Equal(t, expectedJson, actualJson, "Parsing UTML diagram and stringifying the result does not yield the exact same result for file '%s'", path)
	}
}
