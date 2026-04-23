package convert

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/generic"
	pr "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
	"github.com/osingaatje/seshat/types/graph/utml"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

const LINE_CLOSENESS_PERCENT uint8 = 10  // 0-255%, 100% is the full length of the edge
const VERTEX_CLOSENESS_PERCENT uint8 = 2 // 0-255%, 100% = size of node (avg. width and height)

// Parsing raw file to UTML
func parseUTML(c *context.Ctx, filepath string) *ParseResultUTML {
	r, err := os.ReadFile(filepath)
	if err != nil {
		c.LogErr("Error occurred while reading file! Err='%s'", err.Error())
		return nil
	}

	var jsonRes *ParseResultUTML

	err = helper.UnmarshalJSON(r, &jsonRes)
	if err != nil {
		c.LogErr("Could not marshal file '%s' to a UTML Parse Result! Err=%s", filepath, err.Error())
		return nil
	}

	jsonRes.Metadata = NewGraphMetadata(filepath)

	return jsonRes
}

// Converting to internal representation
func convertUTMLToParseRes(c *context.Ctx, utml *ParseResultUTML) *pr.InternalGraph {
	c.LogPrefixAdd("UTML -> Internal '%s'", filepath.Base(utml.Metadata.Filename))
	defer /* I love Go */ c.LogPrefixRm("UTML -> Internal '%s'", filepath.Base(utml.Metadata.Filename))
	if utml == nil {
		c.LogErr("Nil UTML parse result when converting to generic ParseResult.")
		return nil
	}

	res := pr.NewParseResult()
	res.Metadata = utml.Metadata.Copy() // don't forget to add metadata such as filename!

	for i, n := range utml.Nodes {
		vertex, err := convertUTMLVertex(c, i, &n)
		if err != nil { // errors are logged in function
			return nil
		}
		if vertex == nil { // allow for skipping vertices.
			continue
		}

		res.Vertices[vertex.Id] = vertex
	}

	for i, e := range utml.Edges {
		edge := convertUTMLEdge(c, i, &e)
		if edge == nil { // inner errors are logged. no error logging here.
			return nil
		}
		res.Edges[edge.Id] = edge
	}

	// EXTRA METHODS FOR EDGES:
	err := finaliseEdgeProperties(c, utml, res)
	if err != nil {
		c.LogErr(err.Error())
		return nil
	}

	// SANITY CHECKS:
	err = verifyEdgesLinkToVertices(res)
	if err != nil {
		c.LogErr(err.Error())
		return nil
	}

	return res
}

func convertUTMLVertex(c *context.Ctx, index int, n *ParseResultUTMLNode) (*pr.InternalVertex, error) {
	extractedProps := extractUTMLVertexProperties(n)

	if slices.Contains(utml.SKIPPED_VERTEX_TYPES, extractedProps.Type) {
		c.LogDebug("Skipping vertex type '%s', in SKIPPED_TYPES...", extractedProps.Type)
		return nil, nil
	}

	extractedVals := extractUTMLVals(n)
	extractedVisualProps := extractVisualProps(n)

	return &pr.InternalVertex{
		Id:               VertexIdentifier(index), // location in the original UTML array
		Title:            n.Text,
		Properties:       extractedProps,
		Values:           extractedVals,
		VisualProperties: extractedVisualProps,
	}, nil
}
func extractUTMLVertexProperties(n *ParseResultUTMLNode) VertexProperties {
	res := VertexProperties{
		Type:       "",
		Visibility: "",
	}

	res.Type = string(n.Type)
	if n.ClassType != nil {
		res.Type = *n.ClassType
	}
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
		Location:               Vector2D{}.New(n.Position),
		Size:                   Vector2D{X: float64(n.Width), Y: float64(n.Height)},
		VertexStyleFillHex:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // transparent white
		VertexStyleStrokeHex:   color.RGBA{R: 0, G: 0, B: 0, A: 0},         // black
		VertexStyleStrokeWidth: 1,
	}

	if n.HasDoubleBorder {
		res.VertexStyleStrokeWidth = 2 // idk do something fun I guess
	}

	if n.StyleObject != nil {
		res.VertexStyleStrokeWidth = n.StyleObject.StrokeWidth
	}

	return res
}

func convertUTMLEdge(c *context.Ctx, index int, e *ParseResultUTMLEdge) *pr.InternalEdge {
	if e.StartNodeId != nil && *e.StartNodeId < 0 || e.EndNodeId != nil && *e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	res := pr.InternalEdge{
		Id:               EdgeIdentifier(index),
		FromId:           NewVertexIdentifierInt16(e.StartNodeId), // nodeId = location in the array
		ToId:             NewVertexIdentifierInt16(e.EndNodeId),   // nodeId = location in the array
		FromEdgeId:       nil,                                     // filled in later
		ToEdgeId:         nil,                                     // filled in later
		FromProperties:   extractUTMLEdgeEndProps(e, true),
		Label:            extractUTMLEdgeLabel(e),
		ToProperties:     extractUTMLEdgeEndProps(e, false),
		StyleProperties:  extractUTMLEdgeProps(e),
		VisualProperties: EdgeVisualProperties{ /* filled in either with absolute XY in this func, or based on the offset in the finalisation stage */ },
	}

	_tryAddStartEndLocation(c, &res, e.StartPosition, true)
	_tryAddStartEndLocation(c, &res, e.EndPosition, false)

	return &res
}

func _tryAddStartEndLocation(c *context.Ctx, iE *pr.InternalEdge, pos UTMLEdgeXYOrOffsetPosition, isStart bool) bool {
	switch pos := pos.Value.(type) {
	case UTMLEdgeOffsetPosition:
		// nothing happens here, the offset needs to be filled in later based on the Node.

	case UTMLXY:
		if isStart {
			iE.VisualProperties.StartLocation = Vector2D{}.New(pos)
		} else {
			iE.VisualProperties.EndLocation = Vector2D{}.New(pos)
		}
	default:
		c.LogErr("UTML Edge StartPosition was neither an offset nor an X/Y position!")
		return false
	}
	return true
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
		res.Label = &Label{
			Text:     lbl.Value,
			Location: Vector2D{}, // THIS IS DONE AFTER PARSING!
		}
	}

	return res
}

func finaliseEdgeProperties(c *context.Ctx, utml *ParseResultUTML, res *pr.InternalGraph) error {
	for uEID, uE := range utml.Edges {
		nFromId := uE.StartNodeId
		nToId := uE.EndNodeId

		iE, ok := res.Edges[EdgeIdentifier(uEID)]
		if !ok {
			errMsg := fmt.Sprintf("MISSING EDGE - bug in UTML -> Internal conversion. (from,to) = (%d,%d)", uE.StartNodeId, uE.EndNodeId)
			return errors.New(errMsg)
		}

		if nFromId != nil && nToId != nil && (*nFromId < 0 || int(*nFromId) >= len(utml.Nodes) || *nToId < 0 || int(*nToId) >= len(utml.Nodes)) {
			errMsg := fmt.Sprintf("EDGE (%d,%d) HAS INVALID NODE ID(s)", *nFromId, *nToId)
			return errors.New(errMsg)
		}

		var uNodeFrom *ParseResultUTMLNode = nil
		var uNodeTo *ParseResultUTMLNode = nil

		if nFromId != nil {
			uNodeFrom = &utml.Nodes[int(*nFromId)]
		}
		if nToId != nil {
			uNodeTo = &utml.Nodes[int(*nToId)]
		}

		// add start/end location
		err := _addEdgeStartEndLocation(uNodeFrom, uNodeTo, &uE, iE)
		if err != nil {
			c.LogErr("Error while adding start/end location to edge: %s", err.Error())
			return err
		}

		if uE.StartLabel != nil {
			if iE.FromProperties.Label == nil {
				errMsg := fmt.Sprintf("EDGE '%d' HAS BROKEN START LABEL", uEID)
				return errors.New(errMsg)
			}
			_addLocationToLabel(&uE, uE.StartLabel, iE, EdgeLabelPosStart, iE.FromProperties.Label)
		}
		if uE.MiddleLabel != nil {
			if iE.Label == nil {
				errMsg := fmt.Sprintf("EDGE '%d' HAS BROKEN MIDDLE LABEL", uEID)
				return errors.New(errMsg)
			}
			_addLocationToLabel(&uE, uE.MiddleLabel, iE, EdgeLabelPosMiddle, iE.Label)
		}
		if uE.EndLabel != nil {
			if iE.ToProperties.Label == nil {
				errMsg := fmt.Sprintf("EDGE '%d' HAS BROKEN START LABEL", uEID)
				return errors.New(errMsg)
			}
			_addLocationToLabel(&uE, uE.EndLabel, iE, EdgeLabelPosEnd, iE.ToProperties.Label)
		}

	}

	// now that all Edges have a start/end location, we can check if we still need to connect edges to other edges:
	for id, iE := range res.Edges {
		// figure out if we need to add FromEdge / ToEdge
		// this can happen when an edge is connected to only one vertex, or none at all (floating, which is illegal in our case!)
		if iE.FromId == nil || iE.ToId == nil {
			err := tryConnectEdgeEnds(res, id, iE)
			if err != nil {
				return fmt.Errorf("Edge '%d' was not connected to either a starting vertex and/or end vertex - but failed to connect '%d' it to another edge! %s", id, id, err.Error())
			}
		}
	}

	return nil
}

// some edges may not be connected to both nodes. We hope that the creator of the diagram means to connect it to an edge, which we try to implement here.
func tryConnectEdgeEnds(graph *pr.InternalGraph, id EdgeIdentifier, iE *pr.InternalEdge) error {
	if iE.FromId != nil && iE.ToId != nil {
		panic("ONLY USE THIS FUNCTION WHEN A FromId OR ToId IS MISSING")
	}

	// if we've already added edge connections, we're done.
	if ((iE.FromId == nil) != (iE.FromEdgeId == nil)) && ((iE.ToId == nil) != (iE.ToEdgeId == nil)) {
		return nil
	}

	// just so that we don't have to make another function.
	for range 2 {
		var vertexIdToFix **VertexIdentifier
		var edgeIdToFix **EdgeIdentifier
		var location *Vector2D
		var edgeName string

		if iE.FromId == nil && iE.FromEdgeId == nil {
			vertexIdToFix = &iE.FromId
			edgeIdToFix = &iE.FromEdgeId
			location = &iE.VisualProperties.StartLocation
			edgeName = "start"
		} else if iE.ToId == nil && iE.ToEdgeId == nil {
			vertexIdToFix = &iE.ToId
			edgeIdToFix = &iE.ToEdgeId
			location = &iE.VisualProperties.EndLocation
			edgeName = "end"
		} else {
			break
		}

		var smallestDist float64 = math.Inf(1)
		var smallestDistVertexId VertexIdentifier = INVALID_VERT_ID
		var smallestDistEdgeId EdgeIdentifier = INVALID_EDGE_ID

		for id, otherVertex := range graph.Vertices {
			h, w := otherVertex.VisualProperties.Size.X, otherVertex.VisualProperties.Size.Y

			// calculate to distance from the center, then subtract half the size again (dirty hack)
			topLeftPos := otherVertex.VisualProperties.Location
			halfNodeSize := otherVertex.VisualProperties.Size.Mul(0.5)
			centerNodePos := topLeftPos.Add(halfNodeSize)

			halfNodeDistance := topLeftPos.Dist(centerNodePos)
			dist := centerNodePos.Dist(*location)
			dist -= halfNodeDistance

			normalisedDist := dist / ((h + w) / 2)
			if normalisedDist < float64(VERTEX_CLOSENESS_PERCENT) && normalisedDist < smallestDist {
				smallestDist = normalisedDist
				smallestDistVertexId = id
				smallestDistEdgeId = INVALID_EDGE_ID
			}
		}

		for id, otherEdge := range graph.Edges {
			dist, edgeLen := _calculateDistanceToLine(&otherEdge.VisualProperties.StartLocation, &otherEdge.VisualProperties.EndLocation, location)
			if dist/edgeLen < float64(LINE_CLOSENESS_PERCENT) && dist < smallestDist {
				smallestDist = dist
				smallestDistEdgeId = id
				smallestDistVertexId = INVALID_VERT_ID
			}
		}

		if smallestDist < 99999999 {
			if smallestDistVertexId != INVALID_VERT_ID {
				*vertexIdToFix = &smallestDistVertexId
			} else if smallestDistEdgeId != INVALID_EDGE_ID {
				*edgeIdToFix = &smallestDistEdgeId
			} else {
				panic("INVALID CODE: BUG - minimal distance calculation for connecting to a vertex/edge did produce a distance but not an ID")
			}
		} else {
			return fmt.Errorf("Could not find a close enough edge to connect %s point of edge '%d' to located at (%.2f,%.2f)", edgeName, id, location.X, location.Y)
		}
	}
	return nil
}

// see https://mathworld.wolfram.com/Point-LineDistance2-Dimensional.html
func _calculateDistanceToLine(lineStart *Vector2D, lineEnd *Vector2D, loc *Vector2D) (dist float64, edgeLength float64) {
	// calculate distance between the point 'loc(ation)' = (x_0, y_0) and the line lineStart<->lineEnd = (x_1,y_1), (x_2, y_2)
	// we use this formula: d = | v^ . r | = ( |(x_2 - x_1)*(y_1 - y_0) - (x_1 - x_0)(y_2 - y_1)| ) \ ( \sqrt( (x_2-x_1)^2 + (y_2-y_1)^2 ) )

	x_0, y_0 := loc.X, loc.Y
	x_1, y_1 := lineStart.X, lineStart.Y
	x_2, y_2 := lineEnd.X, lineEnd.Y

	lenOfEdge := math.Sqrt(math.Pow(x_2-x_1, 2) + math.Pow(y_2-y_1, 2)) // <>---------------<> the length inbetween lineStart and lineEnd

	topPart := math.Abs((x_2-x_1)*(y_1-y_0) - (x_1-x_0)*(y_2-y_1))
	return topPart / lenOfEdge, lenOfEdge
}

// determine the offset of an edge / label based on edge position
// assumption: we draw nodes/vertices from the top-left.
func _determineXYOffsetBasedOnEdgePos(nodeSize Vector2D, edgePos UTMLEdgeOffsetPosition) Vector2D {
	var res Vector2D
	switch edgePos {
	case EdgePosTopCenter:
		res = Vector2D{X: nodeSize.X / 2, Y: 0}
	case EdgePosTopRight:
		res = Vector2D{X: nodeSize.X, Y: 0}
	case EdgePosMiddleRight:
		res = Vector2D{X: nodeSize.X, Y: nodeSize.Y / 2}
	case EdgePosBottomRight:
		res = nodeSize
	case EdgePosBottomCenter:
		res = Vector2D{X: nodeSize.X / 2, Y: nodeSize.Y}
	case EdgePosBottomLeft:
		res = Vector2D{X: 0, Y: nodeSize.Y}
	case EdgePosMiddleLeft:
		res = Vector2D{X: 0, Y: nodeSize.Y / 2}
	case EdgePosTopLeft:
		res = Vector2D{X: 0, Y: 0}
	}
	return res
}

// Tries to calculate the X/Y start / end position of an edge, if no absolute position is given already.
func _addEdgeStartEndLocation(
	nFrom *ParseResultUTMLNode,
	nTo *ParseResultUTMLNode,
	e *ParseResultUTMLEdge,
	res *pr.InternalEdge) error {
	if e == nil || res == nil {
		panic("you stupid?")
	}

	// figure out if we still need to add the location based on the start/end nodes
	var startIsOffset bool
	var endIsOffset bool

	switch e.StartPosition.Value.(type) {
	case UTMLXY:
		startIsOffset = false
	case UTMLEdgeOffsetPosition:
		startIsOffset = true
	default:
		panic("INVALID START POSITION (not XY/offset")
	}
	switch e.EndPosition.Value.(type) {
	case UTMLXY:
		endIsOffset = false
	case UTMLEdgeOffsetPosition:
		endIsOffset = true
	default:
		panic("INVALID END POSITION (not XY/offset")
	}

	// if we have no nodes to compare to but we get offsets from UTML, we can't do anything... so exit.
	if (nFrom == nil && startIsOffset) || (nTo == nil && endIsOffset) {
		return fmt.Errorf("UTML edge had a start or end offset (not absolute position) but also no Node to reference! Stupid UTML :(")
	}

	// if we have a start offset, calculate it
	if startIsOffset {
		fromPos := Vector2D{}.New(nFrom.Position)
		fromSize := Vector2D{}.NewInt(nFrom.Width, nFrom.Height)

		offsetStart := _determineXYOffsetBasedOnEdgePos(fromSize, e.StartPosition.Value.(UTMLEdgeOffsetPosition) /* cast it to the correct type (is checked beforehand) */)
		res.VisualProperties.StartLocation = fromPos.Add(offsetStart)
	}

	if endIsOffset {
		toPos := Vector2D{}.New(nTo.Position)
		toSize := Vector2D{}.NewInt(nTo.Width, nTo.Height)

		offsetEnd := _determineXYOffsetBasedOnEdgePos(toSize, e.EndPosition.Value.(UTMLEdgeOffsetPosition))
		res.VisualProperties.EndLocation = toPos.Add(offsetEnd)
	}
	return nil
}

// Translates relative position of a UTML label into an absolute position for internal repr.
// ..and adds text etc.
func _addLocationToLabel(
	uE *ParseResultUTMLEdge,
	uL *UTMLEdgeLabel,
	iE *pr.InternalEdge,
	labelPos EdgeLabelPos,
	resL *Label) {
	if uE == nil || uL == nil || resL == nil {
		panic("Internal bug, something is nil.")
	}

	// Edging in progress

	// UTML handles locations very weirdly with a clock-like structure
	// (see UTMLEdgeEndPosition)

	// position depends on
	// - position of label (start/middle/end)
	// - connecting points to the node
	// - offset of that particular edge (if present)
	var position Vector2D
	switch labelPos {
	case EdgeLabelPosStart:
		offset := Vector2D{}.New(uL.Offset)
		// solid base case. We can make it pretty later
		position = iE.VisualProperties.StartLocation.Add(offset)

	case EdgeLabelPosMiddle:
		// the visual properties are to be added before this function. If not, this horribly breaks.
		nFromPos := iE.VisualProperties.StartLocation
		nToPos := iE.VisualProperties.EndLocation

		// if we don't have 'middle positions' then position it between two nodes
		// otherwise take the center-most middle pos.
		if len(uE.MiddlePositions) == 0 {
			position = Vector2D{
				X: (nToPos.X - nFromPos.X) / 2,
				Y: (nToPos.Y - nFromPos.Y) / 2,
			}
		} else {
			var centerPos Vector2D
			middlePosLen := len(uE.MiddlePositions)
			if middlePosLen%2 == 0 {
				rMiddlePos := Vector2D{}.New(uE.MiddlePositions[middlePosLen/2])
				lMiddlePos := Vector2D{}.New(uE.MiddlePositions[middlePosLen/2-1])

				centerPos = Vector2D{
					X: (rMiddlePos.X - lMiddlePos.X) / 2,
					Y: (rMiddlePos.Y - lMiddlePos.Y) / 2,
				}
			} else {
				centerPos = Vector2D{}.New(uE.MiddlePositions[middlePosLen/2])
			}

			offset := Vector2D{}.New(uL.Offset)
			position = centerPos.Add(offset)
		}

	case EdgeLabelPosEnd:
		offset := Vector2D{}.New(uL.Offset)
		// solid base case. We can make it pretty later
		position = iE.VisualProperties.EndLocation.Add(offset)
	}

	resL.Location = position
}

func extractUTMLEdgeProps(e *ParseResultUTMLEdge) EdgeStyleProperties {
	res := EdgeStyleProperties{
		LineStyle: UTMLLineStyleToParsedStyle[e.LineStyle],
		// StartLocation done later
		// EndLocation done later
	}

	return res
}

func extractUTMLEdgeLabel(e *ParseResultUTMLEdge) *Label {
	if e.MiddleLabel == nil {
		return nil
	}

	res := &Label{
		Text: e.MiddleLabel.Value,
		//Location: Vector2D{}.New(e.MiddleLabel.Offset), ADDED LATER!
	}
	return res
}

func verifyEdgesLinkToVertices(r *pr.InternalGraph) error {
	incorrectEdges := []string{}
	for id, e := range r.Edges {
		var hasFromVertex bool = false
		var hasToVertex bool = false
		var hasFromEdge bool = false
		var hasToEdge bool = false

		if e.FromId != nil {
			_, hasFromVertex = r.Vertices[*e.FromId]
		}
		if e.ToId != nil {
			_, hasToVertex = r.Vertices[*e.ToId]
		}
		if e.FromEdgeId != nil {
			_, hasFromEdge = r.Edges[*e.FromEdgeId]
		}
		if e.ToEdgeId != nil {
			_, hasToEdge = r.Edges[*e.ToEdgeId]
		}

		// if the edge has both (or neither) a from vertex/edge or a to vertex/edge, then we did something wrong
		if hasFromVertex == hasFromEdge || hasToVertex == hasToEdge {
			incorrectEdges = append(incorrectEdges, fmt.Sprintf("%d", id))
		}
	}

	if len(incorrectEdges) == 0 {
		return nil
	}

	errMsg := "Some edge(s) was/were not connected properly to either a(n) start/end edge/node! Edge(s): ["
	errMsg += strings.Join(incorrectEdges, ",")
	errMsg += "]"

	return errors.New(errMsg)
}
