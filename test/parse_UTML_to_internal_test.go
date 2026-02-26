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

	const FILEPATH = "./examples/simpleDiag-formatted.utml"

	utml := c.Queries.ParseUTML.Get("Parse UTML", types.ParseUTMLCmd{Filepath: FILEPATH})
	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> internal repr.", utml)

	if utml == nil {
		t.Fatalf("UTML repr. was nil!")
		return
	}

	if intern == nil {
		t.Fatalf("Internal representation was nil!")
		return
	}

	// debugging:
	//utml_json, err := json.Marshal(*utml)
	//if err != nil {
	//	t.Fatalf("could not parse UTML to json. Err=%s", err.Error())
	//	return
	//}
	//intern_json, err := helper.MarshalJSON(intern)
	//if err != nil {
	//	t.Fatalf("Could not parse Internal repr. to JSON. Err=%s", err.Error())
	//}
	// write to file
	//err = os.WriteFile("output_internal_repr.json", []byte(intern_json), os.ModeAppend)
	//if err != nil {
	//	t.Fatalf("failed writing repr. to file, err=%s", err.Error())
	//	return
	//}
	//t.Logf("UTML: \n\t%s", utml_json)
	//t.Logf("INTERNAL: \n\t%s", intern_json)

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

	for i, iVertex := range intern.Vertices {
		if i < 0 || int(i) > len(utml.Nodes) {
			t.Fatalf("A vertex identifier in internal repr. does not match UTML repr.! Id=%d", i)
			return
		}

		uVertex := utml.Nodes[i]

		assert.Equal(t, iVertex.Id, i)
		assert.Equal(t, iVertex.Location.X, uVertex.Position.X)
		assert.Equal(t, iVertex.Location.Y, uVertex.Position.Y)
		if uVertex.ClassType != nil {
			assert.Equal(t, iVertex.Properties.Type, *uVertex.ClassType)
		} else {
			assert.Equal(t, iVertex.Properties.Type, "")
		}

		assert.Equal(t, string(iVertex.Properties.Visibility), "") // utml has no visibility per class.
	}
}
