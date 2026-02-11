package parse

import (
	"encoding/json"
	"os"

	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/context/data"
)

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

func convertUTMLToParseRes(c *context.Ctx, utml *data.ParseResultUTML) *data.ParseResult {
	if utml == nil {
		c.LogErr("Nil UTML parse result when converting to generic ParseResult.")
		return nil
	}

	res := data.NewParseResult()
	for i, n := range utml.Nodes {
		vertex := convertUTMLVertex(i, &n)
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

func convertUTMLVertex(index int, n *data.ParseResultUTMLNode) *data.ParsedVertex {
	titleText := ""
	if n.Text != nil {
		titleText = *n.Text
	}

	return &data.ParsedVertex{
		Id:         uint64(index),
		Title:      titleText,
		Properties: map[data.VertexProperty]string{}, // TODO
		Values:     map[string]data.ParsedValue{},    // TODO
		Location:   data.Location2D{X: n.Position.X, Y: n.Position.Y},
	}
}

func convertUTMLEdge(c *context.Ctx, index int, e *data.ParseResultUTMLEdge) *data.ParsedEdge {
	if e.StartNodeId < 0 || e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := data.ParsedEdge{
		FromId:         uint64(e.StartNodeId),
		ToId:           uint64(e.EndNodeId),
		FromProperties: map[data.EdgeEndProperty]string{}, // TODO
		Label:          data.ParsedLabel{},                // TODO
		ToProperties:   map[data.EdgeEndProperty]string{}, // TODO
	}

	return &res
}
