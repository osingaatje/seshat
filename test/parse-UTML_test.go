package test

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types"
)

func TestUTMLParseSimple(t *testing.T) {
	c := driver.NewContext()

	var FilePaths []string = []string{
		"./examples/simpleDiag-formatted.utml",
		"./examples/multiplicities.utml",
	}

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
		res := c.Queries.ParseUTML.Get("Parse UTML", types.ParseUTMLCmd{Filepath: path})
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
