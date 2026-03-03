package test

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/test/helpers"
	"github.com/osingaatje/seshat/types/command"
	types "github.com/osingaatje/seshat/types/parse-result-utml"
)

func TestConvertSimpleUTMLResultToInternal(t *testing.T) {
	c := driver.NewContext()

	filePaths := helpers.AllUTMLFiles("./examples")

	for _, path := range filePaths {
		verifyRes(c, t, path)
	}
}

func TestConvertBrokenFiles(t *testing.T) {
	filePaths := helpers.AllUTMLFiles("./examples/broken-internal")
	for _, path := range filePaths {
		c := driver.NewContext()

		utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: path})
		assert.NotNil(t, utml) // should only break at the internal conversion

		intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)
		assert.Nil(t, intern)

		assert.True(t, strings.Contains(strings.ToLower(c.Logger.GetLogString()), "could not convert"))
	}
}

func verifyRes(c *context.Ctx, t *testing.T, inputFilePath string) {
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: inputFilePath})
	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)

	if utml == nil {
		t.Fatalf("UTML repr. was nil!")
		return
	}

	if intern == nil {
		t.Fatalf("Internal representation was nil!")
		return
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

		if uEdge.MiddleLabel != nil {
			assert.Equal(t, iEdge.Label.Text, uEdge.MiddleLabel.Value)
		} else {
			assert.Equal(t, "", iEdge.Label.Text)
		}
	}

	for i, iVertex := range intern.Vertices {
		if i < 0 || int(i) > len(utml.Nodes) {
			t.Fatalf("A vertex identifier in internal repr. does not match UTML repr.! Id=%d", i)
			return
		}

		uVertex := utml.Nodes[i]

		assert.Equal(t, iVertex.Id, i)

		// height / width/ location
		assert.Equal(t, iVertex.VisualProperties.Location.X, uVertex.Position.X)
		assert.Equal(t, iVertex.VisualProperties.Location.Y, uVertex.Position.Y)

		assert.Equal(t, iVertex.VisualProperties.Size.X, float64(uVertex.Width))
		assert.Equal(t, iVertex.VisualProperties.Size.Y, float64(uVertex.Height))

		// class type
		if uVertex.ClassType != nil {
			assert.Equal(t, iVertex.Properties.Type, *uVertex.ClassType)
		} else {
			assert.Equal(t, iVertex.Properties.Type, "")
		}

		assert.Equal(t, string(iVertex.Properties.Visibility), "") // utml has no visibility per class.
	}
}
