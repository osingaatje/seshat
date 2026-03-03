package parse

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/parse-result"
	. "github.com/osingaatje/seshat/types/parse-result-utml"
)

// Parsing raw file to UTML
func parseUTML(c *context.Ctx, cmd ParseUTMLCmd) *ParseResultUTML {
	r, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		c.LogErr("Error occurred while reading file! Err='%s'", err.Error())
		return nil
	}

	var jsonRes *ParseResultUTML

	err = helper.UnmarshalJSON(r, &jsonRes)
	if err != nil {
		c.LogErr("Could not marshal file '%s' to a UTML Parse Result! Err=%s", cmd.Filepath, err.Error())
		return nil
	}

	return jsonRes
}

// Converting to internal representation
func convertUTMLToParseRes(c *context.Ctx, utml *ParseResultUTML) *ParseResult {
	const LOGPREFIX = "Could not convert UTML -> internal: "

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

		var eId EdgeIdentifier
		res.Edges[eId.New(edge.FromId, edge.ToId)] = edge
	}

	// EXTRA METHODS FOR EDGES:
	err := finaliseEdgeProperties(utml, res)
	if err != nil {
		c.LogErr("%s%s", LOGPREFIX, err.Error())
		return nil
	}

	// SANITY CHECKS:
	err = verifyEdgesLinkToVertices(res)
	if err != nil {
		c.LogErr("%s%s", LOGPREFIX, err.Error())
		return nil
	}

	return res
}

func convertUTMLVertex(ctx *context.Ctx, index int, n *ParseResultUTMLNode) *ParsedVertex {
	extractedProps := extractUTMLVertexProperties(ctx, index, n)
	extractedVals := extractUTMLVals(n)
	extractedVisualProps := extractVisualProps(n)

	return &ParsedVertex{
		Id:               VertexIdentifier(index), // location in the original UTML array
		Title:            n.Text,
		Properties:       extractedProps,
		Values:           extractedVals,
		VisualProperties: extractedVisualProps,
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

	for _, m := range n.Methods {
		res[m.Name] = ParsedValue{
			Value: "",
			Properties: ValueProperties{
				Visibility: UTMLVisibilityToInternalVisibility[m.Visibility],
				Type:       m.Type,
			},
		}
	}

	return res
}

func extractVisualProps(n *ParseResultUTMLNode) VertexVisualProperties {
	res := VertexVisualProperties{
		Location:               Vector2D{X: n.Position.X, Y: n.Position.Y},
		Size:                   Vector2D{X: float64(n.Width), Y: float64(n.Height)},
		VertexStyleFillHex:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // transparent white
		VertexStyleStrokeHex:   color.RGBA{R: 0, G: 0, B: 0, A: 0},         // black
		VertexStyleStrokeWidth: 1,
	}

	if n.HasDoubleBorder {
		res.VertexStyleStrokeWidth = 2 // idk do something fun I guess
	}

	// ignored, because we don't really need it (now) -2026-03-02
	// if n.StyleObject != nil {
	// 	res.VertexStyleStrokeWidth = n.StyleObject.StrokeWidth
	// }

	return res
}

func convertUTMLEdge(c *context.Ctx, e *ParseResultUTMLEdge) *ParsedEdge {
	if e.StartNodeId < 0 || e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := ParsedEdge{
		FromId:          VertexIdentifier(e.StartNodeId), // location in the array
		ToId:            VertexIdentifier(e.EndNodeId),   // location in the array
		FromProperties:  extractUTMLEdgeEndProps(e, true),
		Label:           extractUTMLEdgeLabel(e),
		ToProperties:    extractUTMLEdgeEndProps(e, false),
		StyleProperties: extractUTMLEdgeProps(e),
	}

	return &res
}

func extractUTMLEdgeEndProps(e *ParseResultUTMLEdge, start bool) EdgeEndProperties {
	res := EdgeEndProperties{
		ArrowStyle: UTMLArrowStyleToInteral[e.StartStyle], // default
		Label:      nil,
	}

	var lbl *UTMLEdgeLabel = e.StartLabel
	if !start {
		lbl = e.EndLabel
	}

	if lbl != nil {
		res.Label = &ParsedLabel{
			Text:     lbl.Value,
			Location: Vector2D{}, // THIS IS DONE AFTER PARSING!
		}
	}

	return res
}

func finaliseEdgeProperties(utml *ParseResultUTML, res *ParseResult) error {
	for _, uE := range utml.Edges {
		nFromId := uE.StartNodeId
		nToId := uE.EndNodeId

		var eId EdgeIdentifier

		iE, ok := res.Edges[eId.New(VertexIdentifier(nFromId), VertexIdentifier(nToId))]
		if !ok {
			errMsg := fmt.Sprintf("MISSING EDGE - bug in UTML -> Internal conversion. (from,to) = (%d,%d)", uE.StartNodeId, uE.EndNodeId)
			return errors.New(errMsg)
		}

		if nFromId < 0 || nToId < 0 || nFromId >= len(utml.Nodes) || nToId >= len(utml.Nodes) {
			errMsg := fmt.Sprintf("EDGE (%d,%d) HAS INVALID NODE ID(s)", nFromId, nToId)
			return errors.New(errMsg)
		}
		uNodeFrom := utml.Nodes[nFromId]
		uNodeTo := utml.Nodes[nFromId]

		if uE.StartLabel != nil {
			if iE.FromProperties.Label == nil {
				errMsg := fmt.Sprintf("EDGE (%d,%d) HAS BROKEN MIDDLE LABEL", uE.StartNodeId, uE.EndNodeId)
				return errors.New(errMsg)
			}
			_addLocationToLabel(&uNodeFrom, &uNodeTo, &uE, uE.StartLabel, iE.FromProperties.Label)
		}
	}

	return nil
}

type EdgeLabelPos int

const (
	EdgeLabelPosStart EdgeLabelPos = iota
	EdgeLabelPosMiddle
	EdgeLabelPosEnd
)

func _addLocationToLabel(
	nFrom *ParseResultUTMLNode,
	nTo *ParseResultUTMLNode,
	e *ParseResultUTMLEdge,
	uL *UTMLEdgeLabel,
	labelPos EdgeLabelPos,
	resL *ParsedLabel) {
	if nFrom == nil || nTo == nil || uL == nil || resL == nil {
		panic("Internal bug, something is nil.")
	}

	// Edging in progress

	// UTML handles locations very weirdly with a clock-like structure
	// (see UTMLEdgeEndPosition)

	// position depends on
	// - position of label (start/middle/end)
	// - connecting points to the node
	// - offset of that particular edge (if present)

	panic("TODO!")
}

func extractUTMLEdgeProps(e *ParseResultUTMLEdge) EdgeStyleProperties {
	res := EdgeStyleProperties{
		LineStyle: UTMLLineStyleToParsedStyle[e.LineStyle],
	}

	return res
}

func extractUTMLEdgeLabel(e *ParseResultUTMLEdge) *ParsedLabel {
	if e.MiddleLabel == nil {
		return nil
	}

	res := &ParsedLabel{
		Text: e.MiddleLabel.Value,
		Location: Vector2D{
			X: e.MiddleLabel.Offset.X,
			Y: e.MiddleLabel.Offset.Y,
		},
	}
	return res
}

func verifyEdgesLinkToVertices(r *ParseResult) error {
	const PREFIX = "Could not convert UTML to internal representation: "

	incorrectEdges := []string{}
	for _, e := range r.Edges {
		_, okFrom := r.Vertices[e.FromId]
		_, okTo := r.Vertices[e.ToId]
		if !okFrom || !okTo {
			incorrectEdges = append(incorrectEdges, fmt.Sprintf("(%d -> %d)", e.FromId, e.ToId))
		}

		if len(incorrectEdges) > 0 {
			errMsg := PREFIX + "Some edge(s) was/were not connected to any nodes! Edge(s): ["
			errMsg += strings.Join(incorrectEdges, ",")
			errMsg += "]"

			return errors.New(errMsg)
		}
	}

	return nil
}
