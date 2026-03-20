package test

import (
	"reflect"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
	"github.com/osingaatje/seshat/types/repair"
)

func TestAllFixableDiagramsShouldBeChangedInSomeWay(t *testing.T) {
	for _, f := range helper.AllUTMLFilesUNSAFE("./examples/fixable") {
		intern, fixed := parseAndFix(t, f)

		assert.NotNil(t, intern)
		assert.NotNil(t, fixed)
		assert.False(t, reflect.DeepEqual(intern, fixed), "Internal and Fixed representations were the same, which contradicts the 'fixable' directory name!")
	}
}

func TestSwapLabelFromTo(t *testing.T) {
	// needs from and to swapped
	intern, fixed := parseAndFix(t, "./examples/fixable/swapped-multiplicites.utml")

	assert.Equal(t, 1, len(intern.Edges))
	assert.Equal(t, len(intern.Edges), len(fixed.Edges))
	assert.Equal(t, len(intern.Vertices), len(fixed.Vertices))

	for id, iE := range intern.Edges {
		fE /*fixed edge */, ok := fixed.Edges[id]
		if !ok {
			t.FailNow()
		}

		assert.NotNil(t, iE.FromProperties.Label)
		assert.NotNil(t, iE.ToProperties.Label)
		assert.NotNil(t, fE.FromProperties.Label)
		assert.NotNil(t, fE.ToProperties.Label)

		assert.Equal(t, *iE.FromProperties.Label, *fE.ToProperties.Label)
		assert.Equal(t, *iE.ToProperties.Label, *fE.FromProperties.Label)
	}
}

func TestSwapMiddleTo(t *testing.T) {
	const FilePath = "./examples/fixable/swap-rightmid.utml"

	intern, fixed := parseAndFix(t, FilePath)

	assert.Equal(t, 1, len(intern.Edges))
	assert.Equal(t, len(intern.Edges), len(fixed.Edges))
	assert.Equal(t, len(intern.Vertices), len(fixed.Vertices))

	for id, iE := range intern.Edges {
		fE /*fixed edge */, ok := fixed.Edges[id]
		if !ok {
			t.FailNow()
		}

		assert.NotNil(t, iE.ToProperties.Label)
		assert.NotNil(t, iE.Label)
		assert.NotNil(t, fE.ToProperties.Label)
		assert.NotNil(t, fE.Label)

		assert.Equal(t, *iE.ToProperties.Label, *fE.Label)
		assert.Equal(t, *iE.Label, *fE.ToProperties.Label)
	}
}

// more complicated case where we need to set one label, and not swap stuff around
func TestReplaceMiddleRight(t *testing.T) {
	// left > right, right > center
	intern, fixed := parseAndFix(t, "./examples/fixable/replace-center-right.utml")

	assert.Equal(t, 1, len(intern.Edges))
	assert.Equal(t, len(intern.Edges), len(fixed.Edges))
	assert.Equal(t, len(intern.Vertices), len(fixed.Vertices))

	for id, iE := range intern.Edges {
		fE /*fixed edge */, ok := fixed.Edges[id]
		if !ok {
			t.FailNow()
		}

		// not-fixed = left & right label (left moved to right, right label moved to center)
		assert.NotNil(t, iE.FromProperties.Label)
		assert.Nil(t, iE.Label)
		assert.NotNil(t, iE.ToProperties.Label)

		assert.Nil(t, fE.FromProperties.Label)
		assert.NotNil(t, fE.Label)
		assert.NotNil(t, fE.ToProperties.Label)

		assert.Equal(t, *iE.FromProperties.Label, *fE.ToProperties.Label)
		assert.Equal(t, *iE.ToProperties.Label, *fE.Label)
	}
}

func parseAndFix(t *testing.T, filePath string) (internal *ParseResult, fixed *ParseResult) {
	c := driver.NewContext()
	utml := c.Queries.ParseUTML.Get("Parse UTML", filePath)
	if utml == nil {
		t.Fatal("Failed to parse UTML")
		return
	}

	internal = c.Queries.ParseUTMLToParseRes.Get("UTML -> internal", utml)
	if internal == nil {
		t.Fatal("Failed to parse internal repr.")
		return
	}

	// check repairs: swap labels from and to
	fixed = c.Queries.RepairDiagram.Get("Repair diagram",
		command.NewRepairCmd(internal, repair.RepairOptions{
			SwapEdgeLabels: true,
		}))
	if fixed == nil {
		t.Fatal("Failed fixing diagram!")
		return
	}
	return internal, fixed
}
