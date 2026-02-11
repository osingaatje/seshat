package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/context/data"
	"github.com/osingaatje/seshat/driver"
)

func TestConvertSimpleUTMLResultToInternal(t *testing.T) {
	c := driver.NewContext()

	const FILEPATH = "./examples/simpleDiag.utml"
	utml := c.Queries.ParseUTML.Get("Parse UTML", data.ParseUTMLCmd{Filepath: FILEPATH})
	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)

	if intern == nil {
		t.Fatalf("Internal representation was nil!")
	}
	assert.Equal(t, len(utml.Edges), len(intern.Edges), "Edges not equal!")
	assert.Equal(t, len(utml.Nodes), len(intern.Vertices), "Nodes not equal!")
	// TODO MORE SHIT
}
