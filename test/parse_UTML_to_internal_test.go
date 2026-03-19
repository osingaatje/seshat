package test

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/src/driver"
	. "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

func TestConvertSimpleUTMLResultToInternal(t *testing.T) {
	c := driver.NewContext()

	filePaths := helper.AllUTMLFilesUNSAFE("./examples/correct")

	for _, path := range filePaths {
		parseAndVerify(c, t, path)
	}
}

func TestConvertBrokenFiles(t *testing.T) {
	filePaths := helper.AllUTMLFilesUNSAFE("./examples/broken-internal")
	for _, path := range filePaths {
		c := driver.NewContext()

		utml := c.Queries.ParseUTML.Get("Parse UTML", path)
		assert.NotNil(t, utml) // should only break at the internal conversion

		intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)
		assert.Nil(t, intern)

		assert.True(t, strings.Contains(strings.ToLower(c.Logger.GetLogString()), "could not convert"))
	}
}

func parseAndVerify(c *context.Ctx, t *testing.T, inputFilePath string) {
	utml := c.Queries.ParseUTML.Get("Parse UTML", inputFilePath)
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
		uEdge, ok := helper.Find(utml.Edges, func(e *ParseResultUTMLEdge) bool {
			return (e.StartNodeId == nil || int64(*e.StartNodeId) == int64(*iEdge.FromId)) &&
				(e.EndNodeId == nil || int64(*e.EndNodeId) == int64(*iEdge.ToId))
		})
		if !ok {
			t.Fatalf("Could not find edge %d -> %d in UTML result", iEdge.FromId, iEdge.ToId)
			return
		}

		// found edge, now compare stuff

		// NODES/VERTICES
		assert.Equal(t, iEdge.FromId == nil, uEdge.StartNodeId == nil) // node connections must be kept
		if iEdge.FromId != nil {
			assert.Equal(t, int64(*iEdge.FromId), int64(*uEdge.StartNodeId))
		}

		assert.Equal(t, iEdge.ToId == nil, uEdge.EndNodeId == nil) // node connections must be kept
		if iEdge.ToId != nil {
			assert.Equal(t, int64(*iEdge.ToId), int64(*uEdge.EndNodeId))
		}
		// END NODES/VERTICES

		// LABELS
		assert.Equal(t, iEdge.FromProperties.Label == nil, uEdge.StartLabel == nil)
		if iEdge.FromProperties.Label != nil {
			assert.Equal(t, iEdge.FromProperties.Label.Text, uEdge.StartLabel.Value)
		}

		assert.Equal(t, iEdge.Label == nil, uEdge.MiddleLabel == nil)
		if iEdge.FromProperties.Label != nil {
			assert.Equal(t, iEdge.FromProperties.Label.Text, uEdge.StartLabel.Value)
		}

		assert.Equal(t, iEdge.ToProperties.Label == nil, uEdge.EndLabel == nil)
		if iEdge.FromProperties.Label != nil {
			assert.Equal(t, iEdge.FromProperties.Label.Text, uEdge.StartLabel.Value)
		}
		// END LABELS

	}

	for i, iVertex := range intern.Vertices {
		if i < 0 || int64(i) > int64(len(utml.Nodes)) {
			t.Fatalf("A vertex identifier in internal repr. does not match UTML repr.! Id=%d", i)
			return
		}

		uVertex := utml.Nodes[i]

		assert.Equal(t, iVertex.Id, i)

		// height / width/ location
		assert.Equal(t, iVertex.VisualProperties.Location, Vector2D{}.New(uVertex.Position))
		assert.Equal(t, iVertex.VisualProperties.Size, Vector2D{}.NewInt(uVertex.Width, uVertex.Height))

		// class type
		if uVertex.ClassType != nil {
			assert.Equal(t, iVertex.Properties.Type, *uVertex.ClassType)
		} else {
			assert.Equal(t, iVertex.Properties.Type, "")
		}

		assert.Equal(t, string(iVertex.Properties.Visibility), "") // utml has no visibility per class.
	}
}
