package test

import (
	"slices"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/src/driver"
	utmlConvert "github.com/osingaatje/seshat/src/queryable/convert/utml"
	"github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/generic"
	"github.com/osingaatje/seshat/types/graph/intern"
	"github.com/osingaatje/seshat/types/graph/utml"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

func TestConvertSpecificFile(t *testing.T) {
	c := driver.NewContext()
	parseAndVerify(c, t, "../DATASETS/2025_M2_TCS/q/5/1027516.json")
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
	filePaths := helper.AllUTMLFilesUNSAFE("./examples/broken-internal")
	for _, path := range filePaths {
		c := driver.NewContext()

		utml := c.Queries.ParseUTML.Get("Parse UTML", path)
		assert.NotNil(t, utml) // should only break at the internal conversion

		intern := c.Queries.ParseUTMLToParseRes.Get("UTML -> internal repr.", utml)
		assert.Nil(t, intern)
	}
}

func parseAndVerify(c *context.Ctx, t *testing.T, inputFilePath string) {
	uGraph := c.Queries.ParseUTML.Get("Parse UTML", inputFilePath)
	iGraph := c.Queries.ParseUTMLToParseRes.Get("UTML -> internal repr.", uGraph)

	if uGraph == nil {
		t.Fatalf("path='%s': UTML repr. was nil!", inputFilePath)
		return
	}

	if iGraph == nil {
		t.Fatalf("path='%s': Internal representation was nil!", inputFilePath)
		return
	}

	verifyConvertedDiag(t, uGraph, iGraph)
	rGraph := c.Queries.RepairDiagram.Get("Repair internal graph", command.NewRepairCmdDefOpt(iGraph))

	// DEBUG
	// dot := c.Queries.DisplayDiagramAsDot.Get("dot", rGraph)
	// c.LogInfo("%s", dot.String())

	verifyRepairedDiag(t, uGraph, rGraph)
}

func verifyConvertedDiag(t *testing.T, utmlGraph *ParseResultUTML, internGraph *intern.InternalGraph) {
	// BASIC CHECKS
	utmlEdgesNotConnectedToSkippedVertices := helper.Filter(utmlGraph.Edges, func(e ParseResultUTMLEdge) bool {
		add := true
		if e.StartNodeId != nil {
			n := utmlGraph.Nodes[(*e.StartNodeId)]
			add = add && !slices.Contains(SKIPPED_VERTEX_TYPES, GetNodeType(&n))
		}
		if e.EndNodeId != nil {
			n := utmlGraph.Nodes[(*e.EndNodeId)]
			add = add && !slices.Contains(SKIPPED_VERTEX_TYPES, GetNodeType(&n))
		}
		return add
	})
	assert.Equal(t, len(utmlEdgesNotConnectedToSkippedVertices), len(internGraph.Edges), "Edges not equal!")

	nonSkippedUTMLNodes := helper.Filter(utmlGraph.Nodes, func(n ParseResultUTMLNode) bool {
		return !slices.Contains(SKIPPED_VERTEX_TYPES, n.Type)
	})
	assert.Equal(t, len(nonSkippedUTMLNodes), len(internGraph.Vertices), "Nodes not equal!")

	// MORE ADVANCED CHECKS
	for id, iEdge := range internGraph.Edges {
		var uEdge *ParseResultUTMLEdge
		var ok bool = int(id) < int(len(utmlGraph.Edges))
		if ok {
			uEdge = &utmlGraph.Edges[int(id)]
		}

		if !ok {
			t.Fatalf("Could not find edge %d -> %d in UTML result", iEdge.FromId, iEdge.ToId)
			return
		}

		// found edge, now compare stuff

		// NODES/VERTICES
		if iEdge.FromId != nil {
			if uEdge.StartNodeId == nil {
				t.Errorf("Adding start nodes to edges should not happen until the 'repair' step!")
			} else {
				assert.Equal(t, int64(*iEdge.FromId), int64(*uEdge.StartNodeId))
			}
		}

		if iEdge.ToId != nil {
			if uEdge.EndNodeId == nil {
				t.Errorf("Adding end nodes to edges should not happen until the 'repair' step!")
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

		assert.Equal(t, iEdge.FromProperties.ArrowStyle, utmlConvert.UTMLArrowStyleToInteral[uEdge.StartStyle], "Start arrow head style should be mapped according to the \"UTMLArrowStyleToInteral\" map!")
		assert.Equal(t, iEdge.ToProperties.ArrowStyle, utmlConvert.UTMLArrowStyleToInteral[uEdge.EndStyle], "End arrow head style should be mapped according to the \"UTMLArrowStyleToInteral\" map!")
		// END LABELS

	}

	for i, iVertex := range internGraph.Vertices {
		if i < 0 || int64(i) > int64(len(utmlGraph.Nodes)) {
			t.Fatalf("A vertex ID in internal repr. does not match UTML repr.! Internal Id=%d", i)
			return
		}

		uVertex := utmlGraph.Nodes[i]

		assert.Equal(t, iVertex.Id, i)

		// height / width/ location
		assert.Equal(t, iVertex.VisualProperties.Location, Vector2D{}.New(uVertex.Position))
		assert.Equal(t, iVertex.VisualProperties.Size, Vector2D{}.NewInt(uVertex.Width, uVertex.Height))

		assert.Equal(t, string(iVertex.Properties.Visibility), "") // utml has no visibility per class.
	}
}

func verifyRepairedDiag(t *testing.T, u *utml.ParseResultUTML, g *intern.InternalGraph) {
	assert.NotNil(t, g, "Repairing a diagram should succeed!")
	// t.Logf("Todo more checks in verify repaired diagram")
}
