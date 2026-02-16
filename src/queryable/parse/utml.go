package parse

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types"
	. "github.com/osingaatje/seshat/types/parse-result-datatypes"
)

// Parsing raw file to UTML
func parseUTML(c *context.Ctx, cmd data.ParseUTMLCmd) *data.ParseResultUTML {
	r, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		c.LogErr("Error occurred while reading file! Err='%s'", err.Error())
		return nil
	}

	var jsonRes *data.ParseResultUTML

	err = json.Unmarshal(r, &jsonRes)
	if err != nil {
		c.LogErr("Could not marshal file '%s' to a UTML Parse Result! Err=%s", cmd.Filepath, err.Error())
		return nil
	}

	return jsonRes
}

// Converting to internal representation
func convertUTMLToParseRes(c *context.Ctx, utml *data.ParseResultUTML) *data.ParseResult {
	if utml == nil {
		c.LogErr("Nil UTML parse result when converting to generic ParseResult.")
		return nil
	}

	res := data.NewParseResult()
	for i, n := range utml.Nodes {
		vertex := convertUTMLVertex(c, i, &n)
		if vertex == nil { // errors are logged in function
			return nil
		}

		res.Vertices[vertex.Id] = vertex
	}

	for i, e := range utml.Edges {
		edge := convertUTMLEdge(c, i, &e)
		if edge == nil { // inner errors are logged. no error logging here.
			return nil
		}
	}

	return res
}

func convertUTMLVertex(ctx *context.Ctx, index int, n *data.ParseResultUTMLNode) *data.ParsedVertex {
	titleText := ""
	if n.Text != nil {
		titleText = *n.Text
	}

	extractedProps, err := extractUTMLVertexProperties(ctx, index, n)
	if err != nil {
		ctx.LogErr("Could not extract properties for node - err=%s", err.Error())
		return nil
	}
	extractedVals, err := map[string]data.ParsedValue{}, nil // TODO!
	ctx.LogDebug("TODO EXTRACT VALS FROM NODE")

	return &data.ParsedVertex{
		Id:         uint64(index), // location in the original UTML array
		Title:      titleText,
		Properties: extractedProps,
		Values:     extractedVals,
		Location:   data.Vector2D{X: n.Position.X, Y: n.Position.Y},
	}
}
func extractUTMLVertexProperties(ctx *context.Ctx, index int, n *data.ParseResultUTMLNode) (map[VertexProperty]any, error) {
	res := map[VertexProperty]any{}

	for prop, typ := range VertexPropertyAll {
		switch prop {

		case VertexPropClassType:
			if n.ClassType == nil {
				ctx.LogWarn("No valid class type for node index '%d'", index)
				continue
			}
			if reflect.TypeOf(*n.ClassType) != typ {
				ctx.LogErr("Data types for VertexPropClassType and UTML ClassType do not match! %T != %T", typ, reflect.TypeOf(*n.ClassType))
				return res, errors.New("Data type mismatch between VertexPropClassType and UTML ClassType")
			}
			res[VertexPropClassType] = n.ClassType
			break

			//case VertexPropClassVisibility:
			//	break
		}
	}

	return res, nil
}

func convertUTMLEdge(c *context.Ctx, index int, e *data.ParseResultUTMLEdge) *data.ParsedEdge {
	if e.StartNodeId < 0 || e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := data.ParsedEdge{
		FromId:         uint64(e.StartNodeId),        // location in the array
		ToId:           uint64(e.EndNodeId),          // location in the array
		FromProperties: map[EdgeEndProperty]string{}, // TODO
		Label:          data.ParsedLabel{},           // TODO
		ToProperties:   map[EdgeEndProperty]string{}, // TODO
	}

	return &res
}
