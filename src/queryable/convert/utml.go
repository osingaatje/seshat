package convert

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"slices"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/generic"
	pr "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/shared"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

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
		vertex, err := convertUTMLVertex(c, i, utml, &n)
		if err != nil { // errors are logged in function
			return nil
		}
		if vertex == nil { // allow for skipping vertices.
			continue
		}

		res.Vertices[vertex.Id] = vertex
	}

	for i, e := range utml.Edges {
		edge, err := convertUTMLEdge(c, i, utml, &e)
		if err != nil { // inner errors are logged. no error logging here.
			return nil
		}
		if edge == nil { // some edges don't need to be converted, so skip them
			continue
		}
		res.Edges[edge.Id] = edge
	}

	// EXTRA METHODS FOR EDGES:
	err := finaliseEdgeProperties(c, utml, res)
	if err != nil {
		c.LogErr(err.Error())
		return nil
	}

	return res
}

func convertUTMLVertex(c *context.Ctx, index int, u *ParseResultUTML, n *ParseResultUTMLNode) (*pr.InternalVertex, error) {
	ntype := GetNodeType(n)
	if slices.Contains(SKIPPED_VERTEX_TYPES, ntype) {
		c.LogDebug("Skipping vertex '%d' because it has skippable type '%s'", index, ntype)
		return nil, nil
	}

	extractedProps := extractUTMLVertexProperties(n)
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

func convertUTMLEdge(c *context.Ctx, index int, u *ParseResultUTML, e *ParseResultUTMLEdge) (*pr.InternalEdge, error) {
	if e.StartNodeId != nil && *e.StartNodeId < 0 || e.EndNodeId != nil && *e.EndNodeId < 0 {
		c.LogErr("FromId or ToId have non-uint64 values for ")
	}

	// remove Edges that are connected to skipped types:
	if e.StartNodeId != nil {
		if (*e.StartNodeId) < 0 || int(*e.StartNodeId) >= len(u.Nodes) {
			return nil, fmt.Errorf("Start node ID '%d' in edge '%d' is not a valid node!", *e.StartNodeId, index)
		}
		// Don't convert edge if it is connected to a Comment Node or other skipped types
		node := u.Nodes[*e.StartNodeId]
		ntype := GetNodeType(&node)
		if slices.Contains(SKIPPED_VERTEX_TYPES, ntype) {
			c.LogDebug("Skipping edge '%d' because its starting node has a skippable type '%s'", index, ntype)
			return nil, nil
		}
	}
	if e.EndNodeId != nil {
		if (*e.EndNodeId) < 0 || int(*e.EndNodeId) >= len(u.Nodes) {
			return nil, fmt.Errorf("Start node ID '%d' in edge '%d' is not a valid node!", *e.EndNodeId, index)
		}
		// Don't convert edge if it is connected to a Comment Node or other skipped types
		node := u.Nodes[*e.EndNodeId]
		ntype := GetNodeType(&node)
		if slices.Contains(SKIPPED_VERTEX_TYPES, ntype) {
			c.LogDebug("Skipping edge '%d' because its starting node has a skippable type '%s'", index, ntype)
			return nil, nil
		}
	}

	// regular conversion stuff
	res := pr.InternalEdge{
		Id:               EdgeIdentifier(index),
		FromId:           NewVertexIdentifierInt(e.StartNodeId), // nodeId = location in the array
		ToId:             NewVertexIdentifierInt(e.EndNodeId),   // nodeId = location in the array
		FromEdgeId:       nil,                                   // filled in later
		ToEdgeId:         nil,                                   // filled in later
		FromProperties:   extractUTMLEdgeEndProps(e, true),
		Label:            extractUTMLEdgeLabel(e),
		ToProperties:     extractUTMLEdgeEndProps(e, false),
		StyleProperties:  extractUTMLEdgeProps(e),
		VisualProperties: EdgeVisualProperties{ /* filled in either with absolute XY in this func, or based on the offset in the finalisation stage */ },
	}

	err := _tryAddPath(c, &res, e)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func _tryAddPath(c *context.Ctx, iE *pr.InternalEdge, e *ParseResultUTMLEdge) error {
	path := []Vector2D{{}, {}} //always start with start/end

	for isStart, pos := range map[bool]UTMLEdgeXYOrOffsetPosition{
		true:  e.StartPosition,
		false: e.EndPosition,
	} {
		switch pos := pos.Value.(type) {
		case UTMLEdgeOffsetPosition:
			// nothing happens here, the offset needs to be filled in later based on the Node.

		case UTMLXY:
			if isStart {
				path[0] = Vector2D{}.New(pos)
			} else {
				path[1] = Vector2D{}.New(pos)
			}
		default:
			return fmt.Errorf("UTML Edge StartPosition was neither an offset nor an X/Y position!")
		}
	}

	if len(e.MiddlePositions) > 0 {
		last := path[1]
		path = path[:1] // only start

		for _, mpos := range e.MiddlePositions {
			path = append(path, Vector2D{}.New(mpos))
		}
		path = append(path, last)
	}
	iE.VisualProperties.Path = path

	return nil
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
			// apparently this edge was skipped. Ignore it.
			continue
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

	return nil
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
	if len(res.VisualProperties.Path) < 2 {
		panic("INTERNAL EDGE SHOULD ALWAYS HAVE AT LEAST TWO PATH LOCATIONS!!")
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
		res.VisualProperties.Path[0] = fromPos.Add(offsetStart)
	}

	if endIsOffset {
		toPos := Vector2D{}.New(nTo.Position)
		toSize := Vector2D{}.NewInt(nTo.Width, nTo.Height)

		offsetEnd := _determineXYOffsetBasedOnEdgePos(toSize, e.EndPosition.Value.(UTMLEdgeOffsetPosition))
		res.VisualProperties.Path[len(res.VisualProperties.Path)-1] = toPos.Add(offsetEnd)
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
	if len(iE.VisualProperties.Path) < 2 {
		panic("INTERNAL EDGE MUST HAVE >= 2 PATH POSITIONS!")
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
		position = iE.VisualProperties.Path[0].Add(offset)

	case EdgeLabelPosMiddle:
		centerPos := helper.GetCenterPos(iE.VisualProperties.Path)
		offset := Vector2D{}.New(uL.Offset)
		position = centerPos.Add(offset)

	case EdgeLabelPosEnd:
		offset := Vector2D{}.New(uL.Offset)
		// solid base case. We can make it pretty later
		position = iE.VisualProperties.Path[len(iE.VisualProperties.Path)-1].Add(offset)
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
