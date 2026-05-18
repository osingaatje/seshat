package display

import (
	"fmt"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/graph/dot"
	. "github.com/osingaatje/seshat/types/graph/intern"
)

const POSITION_SCALE_FACTOR float64 = float64(1) / 40

func FindQueries(c *context.Ctx) {
	c.Queries.DisplayDiagramAsDot = context.DefineQuery(
		c,
		"Display internal graph representation as a .dot file",
		convertInternalToDot,
	)
}

func convertInternalToDot(c *context.Ctx, p *InternalGraph) *DotGraph {
	if p == nil {
		c.LogWarn("Nil ParseResult in call to Dot conversion!")
		return nil
	}

	res := DotGraph{}.New("parse_result")
	res.NodeSettings = DotNodeSettings{
		Shape: "rect",
	}

	for _, v := range p.Vertices {
		res.Nodes = append(res.Nodes, convertNode(v))
	}
	for _, e := range p.Edges {
		res.Edges = append(res.Edges, convertEdge(p, e))
	}

	return &res
}

func convertNode(v *InternalVertex) DotNode {
	res := DotNode{
		Text:     getNodeText(v),
		NodeOpts: extractNodeOpts(v),
	}

	return res
}
func getNodeText(v *InternalVertex) string {
	res := fmt.Sprintf("________ (id '%d') %s ________", v.Id, v.Title)

	vals := []string{}

	for key, val := range v.Values {
		value := val.Value
		if val.Value != "" {
			value = " default='" + val.Value + "'"
		}
		vals = append(vals, fmt.Sprintf("%s%s : %s (%s)", key, value, val.Properties.Type, val.Properties.Visibility))
	}
	slices.Sort(vals) // otherwise the element ordering is random and graphviz doesn't recognise it as the same node

	if len(vals) > 0 {
		res += "\n" + strings.Join(vals, "\n")
	}

	return res
}

func getEdgeIdText(e *InternalEdge) string {
	return fmt.Sprintf("EID '%d'", e.Id)
}

func extractNodeOpts(v *InternalVertex) DotNodeOptions {
	newVec := v.VisualProperties.Location.Mul(POSITION_SCALE_FACTOR).MulComponents(1, -1) // Y SCALE IS INVERTED IN DOT!
	res := DotNodeOptions{
		Pos:       &newVec,
		NoJustify: true,
	}
	return res
}

func convertEdge(r *InternalGraph, e *InternalEdge) DotEdge {
	res := DotEdge{
		IdText:      getEdgeIdText(e),
		FromText:    "",
		MiddleLabel: "",
		ToText:      "",

		EdgeOpts: DotEdgeOptions{
			Style: EdgeStyleToGraphvizStyle[e.StyleProperties.LineStyle],
		},
	}

	if e.FromId != nil {
		res.FromText = getNodeText(r.Vertices[*e.FromId])
	}
	if e.FromEdgeId != nil {
		res.FromText = getEdgeIdText(r.Edges[*e.FromEdgeId])
	}

	if e.ToId != nil {
		res.ToText = getNodeText(r.Vertices[*e.ToId])
	}
	if e.ToEdgeId != nil {
		res.ToText = getEdgeIdText(r.Edges[*e.ToEdgeId])
	}

	if e.FromProperties.Label != nil {
		res.StartLabel = &e.FromProperties.Label.Text
	}
	if e.Label != nil {
		res.MiddleLabel = e.Label.Text
	}
	if e.ToProperties.Label != nil {
		res.EndLabel = &e.ToProperties.Label.Text
	}

	return res
}
