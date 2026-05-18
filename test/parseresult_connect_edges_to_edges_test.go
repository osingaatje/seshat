package test

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types/command"
	pr "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

func TestConnectEdgeToOtherEdge(t *testing.T) {
	c := driver.NewContext()

	// this file has an association class-type structure, with a
	// dotted line to another edge.
	// This should be converted in the parse result to an Edge connection
	// ..which encodes the semantics of the association class.
	utml := c.Queries.ParseUTML.Get("Parse UTML", "./examples/correct/association-class-simplified.utml")
	if utml == nil {
		t.Fatal("Failed to parse UTML")
		return
	}

	internal := c.Queries.ParseUTMLToParseRes.Get("UTML -> internal repr.", utml)
	if internal == nil {
		t.Fatal("Failed to parse internal repr.")
		return
	}
	repairRes := c.Queries.RepairDiagram.Get("Fix internal rep.", command.NewRepairCmdDefOpt(internal))
	if len(repairRes.Errors) > 0 || repairRes.Diagram == nil {
		t.Fatal("Failed to repair internal repr.")
		return
	}
	internal = repairRes.Diagram

	//debug conversion to .dot file:
	// dot := c.Queries.DisplayDiagramAsDot.Get("dot", internal)
	// if dot == nil {
	// 	panic("DAADLSKFSDLKFJSD")
	// }
	// t.Log(dot.String())

	// vertex names: "OneClass" 1-----* "ManyClass" ---- "ShouldNotConnect"
	// 								|
	// 								| <--- dotted
	// 						"AssociationClass"

	oneClsId, _, _ := _findVertexByNameOrFail(t, internal, "OneClass")
	manyClsId, _, _ := _findVertexByNameOrFail(t, internal, "ManyClass")
	assClsId, _, _ := _findVertexByNameOrFail(t, internal, "AssociationClass")

	oneToManyEdgeId, _, ok := _findEdgeByFunc(t, internal, func(e *pr.InternalEdge) bool {
		return e.FromId != nil && (*e.FromId) == (*oneClsId) && e.ToId != nil && (*e.ToId) == (*manyClsId)
	})
	if !ok {
		t.Fatal("Could not find 1..* edge")
		return
	}

	_, assEdge, ok := _findEdgeByFunc(t, internal, func(e *pr.InternalEdge) bool {
		return e.FromId == nil && e.ToId != nil && (*e.ToId) == (*assClsId)
	})
	if !ok {
		t.Fatal("Could not find association edge")
	}

	assert.Nil(t, assEdge.FromId)
	assert.NotNil(t, assEdge.FromEdgeId)
	if assEdge.FromEdgeId == nil {
		return
	}
	assert.Equal(t, *assEdge.FromEdgeId, *oneToManyEdgeId)
}

func _findVertexByNameOrFail(t *testing.T, graph *pr.InternalGraph, vertexName string) (*VertexIdentifier, *pr.InternalVertex, bool) {
	id, pv, ok := helper.FindValue(graph.Vertices, func(pv *pr.InternalVertex) bool {
		return pv.Title == vertexName
	})

	if !ok {
		t.Fatalf("Could not find vertex named '%s'", vertexName)
		return nil, nil, false
	}
	return id, *pv, ok
}

func _findEdgeByFunc(t *testing.T, graph *pr.InternalGraph, f func(e *pr.InternalEdge) bool) (*EdgeIdentifier, *pr.InternalEdge, bool) {
	id, pe, ok := helper.FindValue(graph.Edges, f)
	if !ok {
		return nil, nil, false
	}
	return id, *pe, true
}
