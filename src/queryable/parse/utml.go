package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types"
	. "github.com/osingaatje/seshat/types/parse-result-datatypes"
)

// Parsing raw file to UTML
func parseUTML(c *context.Ctx, cmd ParseUTMLCmd) *ParseResultUTML {
	r, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		c.LogErr("Error occurred while reading file! Err='%s'", err.Error())
		return nil
	}

	var jsonRes *ParseResultUTML

	err = json.Unmarshal(r, &jsonRes)
	if err != nil {
		c.LogErr("Could not marshal file '%s' to a UTML Parse Result! Err=%s", cmd.Filepath, err.Error())
		return nil
	}

	return jsonRes
}

// Converting to internal representation
func convertUTMLToParseRes(c *context.Ctx, utml *ParseResultUTML) *ParseResult {
	if utml == nil {
		c.LogErr("Nil UTML parse result when converting to generic ParseResult.")
		return nil
	}

	res := NewParseResult()
	for i, n := range utml.Nodes {
		vertex := convertUTMLVertex(c, i, &n)
		if vertex == nil { // errors are logged in function
			return nil
		}

		res.Vertices[VertexIdentifier(vertex.Id)] = vertex
	}

	for _, e := range utml.Edges {
		edge := convertUTMLEdge(c, &e)
		if edge == nil { // inner errors are logged. no error logging here.
			return nil
		}
		res.Edges[EdgeIdentifier{FromId: VertexIdentifier(edge.FromId), ToId: VertexIdentifier(edge.ToId)}] = edge
	}

	// SANITY CHECKS:
	verifyEdgesLinkToVertices(c, res)

	return res
}

func convertUTMLVertex(ctx *context.Ctx, index int, n *ParseResultUTMLNode) *ParsedVertex {
	titleText := ""
	if n.Text != nil {
		titleText = *n.Text
	}

	extractedProps, err := extractUTMLVertexProperties(ctx, index, n)
	if err != nil {
		ctx.LogErr("Could not extract properties for node - err=%s", err.Error())
		return nil
	}
	extractedVals, err := map[string]ParsedValue{}, nil // TODO!
	ctx.LogDebug("TODO EXTRACT VALS FROM NODE")

	return &ParsedVertex{
		Id:         VertexIdentifier(index), // location in the original UTML array
		Title:      titleText,
		Properties: extractedProps,
		Values:     extractedVals,
		Location:   Vector2D{X: n.Position.X, Y: n.Position.Y},
	}
}
func extractUTMLVertexProperties(ctx *context.Ctx, index int, n *ParseResultUTMLNode) (map[VertexProperty]any, error) {
	res := map[VertexProperty]any{}

	for prop, typ := range VertexPropertyAll {
		switch prop {

		case VertexPropClassType:
			if n.ClassType == nil {
				ctx.LogWarn("No valid class type for node index '%d'", index)
				continue
			}
			if reflect.TypeOf(*n.ClassType) != typ {
				ctx.LogErr("types for VertexPropClassType and UTML ClassType do not match! %T != %T", typ, reflect.TypeOf(*n.ClassType))
				return res, errors.New("type mismatch between VertexPropClassType and UTML ClassType")
			}
			res[VertexPropClassType] = n.ClassType

		case VertexPropClassVisibility:
			// nodes do not have visibility :(
		}
	}

	return res, nil
}

func convertUTMLEdge(c *context.Ctx, e *ParseResultUTMLEdge) *ParsedEdge {
	if e.StartNodeId < 0 || e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := ParsedEdge{
		FromId:          VertexIdentifier(e.StartNodeId), // location in the array
		ToId:            VertexIdentifier(e.EndNodeId),   // location in the array
		FromProperties:  extractUTMLEdgeFromProps(c, e),
		Label:           ParsedLabel{}, // TODO
		ToProperties:    extractUTMLEdgeToProps(c, e),
		StyleProperties: extractUTMLEdgeProps(c, e),
	}

	return &res
}

func extractUTMLEdgeFromProps(c *context.Ctx, e *ParseResultUTMLEdge) map[EdgeEndProperty]any {
	res := map[EdgeEndProperty]any{}

	for prop, _ := range EdgeEndPropertyAll {
		switch prop {
		case EdgeEndPropArrowStyle:
			// e.StartStyle
			c.LogErr("TODO ARROW STYLE")
		case EdgeEndPropMultiplicity: // inspect the code for a VertexStartStyle that looks like "*", "0..1", "0..*" etc.
			c.LogErr("TODO PROP MULTIPLICITY")
		}
	}
	return res
}

func extractUTMLEdgeToProps(c *context.Ctx, e *ParseResultUTMLEdge) map[EdgeEndProperty]any {
	res := map[EdgeEndProperty]any{}

	for prop := range EdgeEndPropertyAll {
		switch prop {
		case EdgeEndPropArrowStyle:
			c.LogErr("TODO ARROW STYLE")
		case EdgeEndPropMultiplicity: // inspect the code for a VertexStartStyle that looks like "*", "0..1", "0..*" etc.
			c.LogErr("TODO MULTIPLICITY")
		}
	}
	return res
}

func extractUTMLEdgeProps(c *context.Ctx, e *ParseResultUTMLEdge) map[EdgeStyleProperty]any {
	res := map[EdgeStyleProperty]any{}

	for prop := range EdgeStylePropertyAll {
		switch prop {
		case EdgeStyleLine:
			if e.LineStyle == nil {
				continue
			}

			lineStyle, ok := UTMLLineStyleToParsedStyle[*e.LineStyle]
			if !ok {
				c.LogErr("Unknown UTML line style '%d'! Please make a translation for this", *e.LineStyle)
				continue
			}

			res[EdgeStyleLine] = lineStyle
		}
	}
	return res
}

func verifyEdgesLinkToVertices(c *context.Ctx, r *ParseResult) {
	incorrectEdges := []string{}
	for _, e := range r.Edges {
		_, okFrom := r.Vertices[e.FromId]
		_, okTo := r.Vertices[e.ToId]
		if !okFrom || !okTo {
			incorrectEdges = append(incorrectEdges, fmt.Sprintf("(%d -> %d)", e.FromId, e.ToId))
		}

		if len(incorrectEdges) > 0 {
			errMsg := "Some edges were not connected to any nodes! Edges: "
			errMsg += strings.Join(incorrectEdges, ",")
			c.LogErr(errMsg)
		}
	}
}
