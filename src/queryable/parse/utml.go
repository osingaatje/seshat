package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/osingaatje/seshat/helper"
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

	extractedProps := extractUTMLVertexProperties(ctx, index, n)
	extractedVals := extractUTMLVals(n)

	return &ParsedVertex{
		Id:         VertexIdentifier(index), // location in the original UTML array
		Title:      titleText,
		Properties: extractedProps,
		Values:     extractedVals,
		Location:   Vector2D{X: n.Position.X, Y: n.Position.Y},
	}
}
func extractUTMLVertexProperties(ctx *context.Ctx, index int, n *ParseResultUTMLNode) VertexProperties {
	res := VertexProperties{
		Type:       "",
		Visibility: "",
	}

	// Type:
	if n.ClassType == nil {
		ctx.LogWarn("No valid class type for node index '%d'", index)
		return res
	}

	res.Type = *n.ClassType
	// Visibility: not present for nodes in UTML :(

	return res
}

func extractUTMLVals(n *ParseResultUTMLNode) map[string]ParsedValue {
	res := map[string]ParsedValue{}

	for _, a := range n.Attributes {
		res[a.Name] = ParsedValue{
			Value: "", // utml cannot have (default) values
			Properties: ValueProperties{
				Visibility: UTMLVisibilityToInternalVisibility[a.Visibility],
				Type:       a.Type,
			},
		}
	}

	return res
}

func convertUTMLEdge(c *context.Ctx, e *ParseResultUTMLEdge) *ParsedEdge {
	if e.StartNodeId < 0 || e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := ParsedEdge{
		FromId:          VertexIdentifier(e.StartNodeId), // location in the array
		ToId:            VertexIdentifier(e.EndNodeId),   // location in the array
		FromProperties:  extractUTMLEdgeEndProps(c, e, true),
		Label:           extractUTMLEdgeLabel(e),
		ToProperties:    extractUTMLEdgeEndProps(c, e, false),
		StyleProperties: extractUTMLEdgeProps(e),
	}

	return &res
}

func extractUTMLEdgeEndProps(c *context.Ctx, e *ParseResultUTMLEdge, start bool) EdgeEndProperties {
	res := EdgeEndProperties{
		ArrowStyle:   ArrowStyleNoArrow, // default
		Multiplicity: nil,
	}

	var style *UTMLArrowHeadStyle = e.StartStyle
	var lbl *UTMLEdgeLabel = e.StartLabel
	if !start {
		style = e.EndStyle
		lbl = e.EndLabel
	}

	if style != nil {
		res.ArrowStyle = UTMLArrowStyleToInteral[*style]
	}

	if lbl != nil {
		mult, ok := helper.GetMultiplicity(lbl.Value)
		if !ok {
			c.LogDebug("Could not parse UTML StartLabel ('from') '%s' into a multiplicity.", e.StartLabel.Value)
		}
		res.Multiplicity = mult
	}

	return res
}

func extractUTMLEdgeProps(e *ParseResultUTMLEdge) EdgeStyleProperties {
	res := EdgeStyleProperties{
		LineStyle: EdgeLineStyleSolid,
	}

	if e.LineStyle != nil {
		res.LineStyle = UTMLLineStyleToParsedStyle[*e.LineStyle]
	}

	return res
}

func extractUTMLEdgeLabel(e *ParseResultUTMLEdge) ParsedLabel {
	res := ParsedLabel{Text: "", Location: Vector2D{X: 0, Y: 0}}

	if e.MiddleLabel != nil {
		res.Text = e.MiddleLabel.Value
		res.Location.X = e.MiddleLabel.Offset.X
		res.Location.Y = e.MiddleLabel.Offset.Y
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
