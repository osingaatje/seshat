package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/src/driver"
	. "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

func TestConvertSpecificFile(t *testing.T) {
	c := driver.NewContext()
	parseAndVerify(c, t, "../DATASETS/2025_M2_BIT/q/1/121464.json")
}

func TestConvertAllDatasetFiles(t *testing.T) {
	c := driver.NewContext()

	filePaths, err := helper.AllDatasetFiles()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	for _, path := range filePaths {
		parseAndVerify(c, t, path)
	}
}

func TestConvertSimpleUTMLResultToInternal(t *testing.T) {
	c := driver.NewContext()

	filePaths := helper.AllUTMLFilesUNSAFE("./examples/correct")

	for _, path := range filePaths {
		parseAndVerify(c, t, path)
	}
}

func TestConvertBrokenFiles(t *testing.T) {
	filePaths := helper.AllUTMLFilesUNSAFE("./examples/broken-intern")
	for _, path := range filePaths {
		c := driver.NewContext()

		utml := c.Queries.ParseUTML.Get("Parse UTML", path)
		assert.NotNil(t, utml) // should only break at the internal conversion

		intern := c.Queries.ParseUTMLToParseRes.Get("UTML -> internal repr.", utml)
		assert.Nil(t, intern)
	}
}

func parseAndVerify(c *context.Ctx, t *testing.T, inputFilePath string) {
	utml := c.Queries.ParseUTML.Get("Parse UTML", inputFilePath)
	intern := c.Queries.ParseUTMLToParseRes.Get("UTML -> internal repr.", utml)

	if utml == nil {
		t.Fatalf("path='%s': UTML repr. was nil!", inputFilePath)
		return
	}

	if intern == nil {
		t.Fatalf("path='%s': Internal representation was nil!", inputFilePath)
		return
	}

	// BASIC CHECKS
	assert.Equal(t, len(utml.Edges), len(intern.Edges), "Edges not equal!")
	assert.Equal(t, len(utml.Nodes), len(intern.Vertices), "Nodes not equal!")

	// MORE ADVANCED CHECKS
	for id, iEdge := range intern.Edges {
		var uEdge *ParseResultUTMLEdge
		var ok bool = int(id) < int(len(utml.Edges))
		if ok {
			uEdge = &utml.Edges[int(id)]
		}
		/* OLD COMPLICATED LOGIC: helper.Find(utml.Edges, func(e *ParseResultUTMLEdge) bool {

		// match based on start/end location if there is no start/end node
		match := true
		if e.StartNodeId != nil {
			match = match && iEdge.FromId != nil && int64(*iEdge.FromId) == int64(*e.StartNodeId)
		} else {
			match = match && iEdge.FromId == nil
			switch val := e.StartPosition.Value.(type) {
			case UTMLXY:
				match = match && Vector2D{}.New(val) == iEdge.VisualProperties.StartLocation
			default:
				return false
			}
		}

		if e.EndNodeId != nil {
			match = match && iEdge.ToId != nil && int64(*iEdge.ToId) == int64(*e.EndNodeId)
		} else {
			match = match && iEdge.FromId == nil
			switch val := e.EndPosition.Value.(type) {
			case UTMLXY:
				match = match && Vector2D{}.New(val) == iEdge.VisualProperties.EndLocation
			default:
				return false
			}
		}

		return match
		})*/

		if !ok {
			t.Fatalf("Could not find edge %d -> %d in UTML result", iEdge.FromId, iEdge.ToId)
			return
		}

		// found edge, now compare stuff

		// NODES/VERTICES
		if iEdge.FromId != nil {
			if uEdge.StartNodeId == nil {
				t.Logf("TODO VERIFY WHETHER THE NODE IS WITHIN DISTANCE OF THE EDGE")
			} else {
				assert.Equal(t, int64(*iEdge.FromId), int64(*uEdge.StartNodeId))
			}
		}

		if iEdge.ToId != nil {
			if uEdge.EndNodeId == nil {
				t.Logf("TODO VERIFY WHETHER THE NODE IS WITHIN DISTANCE OF THE EDGE")
			} else {
				assert.Equal(t, int64(*iEdge.ToId), int64(*uEdge.EndNodeId))
			}
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

		assert.Equal(t, string(iVertex.Properties.Visibility), "") // utml has no visibility per class.
	}
}
