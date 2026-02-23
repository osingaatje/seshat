package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types"
)

func TestConvertSimpleUTMLResultToInternal(t *testing.T) {
	c := driver.NewContext()

	const FILEPATH = "./examples/simpleDiag.utml"
	utml := c.Queries.ParseUTML.Get("Parse UTML", types.ParseUTMLCmd{Filepath: FILEPATH})
	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)

	if intern == nil {
		t.Fatalf("Internal representation was nil!")
	}

	// BASIC CHECKS
	assert.Equal(t, len(utml.Edges), len(intern.Edges), "Edges not equal!")
	assert.Equal(t, len(utml.Nodes), len(intern.Vertices), "Nodes not equal!")

	// MORE ADVANCED CHECKS
	for _, iEdge := range intern.Edges {
		uEdge, ok := helper.Find(utml.Edges, func(e *types.ParseResultUTMLEdge) bool {
			return e.StartNodeId == int(iEdge.FromId) && e.EndNodeId == int(iEdge.ToId)
		})
		if !ok {
			t.Fatalf("Could not find edge %d -> %d in UTML result", iEdge.FromId, iEdge.ToId)
			return
		}

		// found edge, now compare shit
		assert.Equal(t, int(iEdge.FromId), uEdge.StartNodeId)
		assert.Equal(t, int(iEdge.ToId), uEdge.EndNodeId)
		assert.Equal(t, iEdge.Label.Text, uEdge.MiddleLabel.Value)
	}
}
