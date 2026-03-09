package test

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/repair"
)

func TestSwapLabelFromTo(t *testing.T) {
	// needs from and to swapped
	const FilePath = "./examples/fixable/swapped-multiplicites.utml"

	c := driver.NewContext()
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: FilePath})
	if utml == nil {
		t.Fatal("Failed to parse UTML")
		return
	}

	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal", utml)
	if intern == nil {
		t.Fatal("Failed to parse internal repr.")
		return
	}
	byt, err := helper.MarshalJSONWithIndent(intern)
	if err != nil {
		t.FailNow()
		return
	}
	os.WriteFile("./EXPORT.json", byt, os.ModePerm)

	// check repairs: swap labels from and to
	fixed := c.Queries.RepairDiagram.Get("Repair diagram",
		command.NewRepairCmd(intern, repair.RepairOptions{
			SwapEdgeLabels: true,
		}))
	if fixed == nil {
		t.Fatal("Failed fixing diagram!")
		return
	}

	assert.Equal(t, 1, len(intern.Edges))
	assert.Equal(t, len(intern.Edges), len(fixed.Edges))
	assert.Equal(t, len(intern.Vertices), len(fixed.Vertices))

	assert.NotNil(t, intern.Edges[0].FromProperties.Label)
	assert.NotNil(t, intern.Edges[0].ToProperties.Label)
	assert.NotNil(t, fixed.Edges[0].FromProperties.Label)
	assert.NotNil(t, fixed.Edges[0].ToProperties.Label)

	assert.Equal(t, *intern.Edges[0].FromProperties.Label, *fixed.Edges[0].ToProperties.Label)
	assert.Equal(t, *intern.Edges[0].ToProperties.Label, *fixed.Edges[0].FromProperties.Label)
}
