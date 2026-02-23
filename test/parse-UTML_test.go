package test

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/test/helpers"
	"github.com/osingaatje/seshat/types"
)

func TestUTMLParseSimple(t *testing.T) {
	c := driver.NewContext()

	const FILEPATH = "./examples/simpleDiag.utml"

	fileContents, err := os.ReadFile(FILEPATH)
	if err != nil {
		t.Errorf("Could not read file: Err=%s", err.Error())
		return
	}

	expectedJson, err := helpers.IndentJSON(fileContents)
	if err != nil {
		t.Errorf("Could not indent JSON for file '%s', err=%s", FILEPATH, err.Error())
		return
	}

	// run query
	res := c.Queries.ParseUTML.Get("Parse UTML", types.ParseUTMLCmd{Filepath: FILEPATH})
	actualJson, err := helpers.IndentJSON([]byte(res.String()))
	if err != nil {
		t.Errorf("Could not indent JSON for parsed UTML, err=%s", err.Error())
		return
	}

	// Write to temp file - useful to diff the JSON in an actual editor.
	//	err = os.WriteFile("test.utml", []byte(actualJson), os.ModePerm)
	//	if err != nil {
	//		t.FailNow()
	//	}

	assert.Equal(t, expectedJson, actualJson, "Parsing UTML diagram and stringifying the result does not yield the exact same result for file '%s'", FILEPATH)
}
