package display

import (
	"fmt"
	"slices"
	"strings"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/graph/dot"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
)

const POSITION_SCALE_FACTOR float64 = float64(1) / 8

func FindQueries(c *context.Ctx) {
	c.Queries.DisplayDiagramAsDot = context.DefineQuery(
		c,
		"Display internal graph representation as a .dot file",
		convertInternalToDot,
	)
}

func convertInternalToDot(c *context.Ctx, p *ParseResult) *DotGraph {
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

func convertNode(v *ParsedVertex) DotNode {
	res := DotNode{
		Text:     extractTextFromNode(v),
		NodeOpts: extractNodeOpts(v),
	}

	return res
}
func extractTextFromNode(v *ParsedVertex) string {
	res := v.Title
	vals := []string{}
	for key, val := range v.Values {
		value := val.Value
		if val.Value != "" {
			value = " default='" + val.Value + "'"
		}
		vals = append(vals, fmt.Sprintf("%s%s type='%s',vis='%s'", key, value, val.Properties.Type, val.Properties.Visibility))
	}
	slices.Sort(vals) // otherwise the element ordering is random and graphviz doesn't recognise it as the same node

	if len(vals) > 0 {
		res += "\nvals='" + strings.Join(vals, ",") + "'"
	}

	res = fmt.Sprintf("\"%s\"", res)

	return res
}

func extractNodeOpts(v *ParsedVertex) DotNodeOptions {
	newVec := v.VisualProperties.Location.Mul(POSITION_SCALE_FACTOR)
	res := DotNodeOptions{
		Pos:       &newVec,
		NoJustify: true,
	}
	return res
}

func convertEdge(r *ParseResult, e *ParsedEdge) DotEdge {
	res := DotEdge{
		FromText: extractTextFromNode(r.Vertices[e.FromId]),
		ToText:   extractTextFromNode(r.Vertices[e.ToId]),
	}

	if e.FromProperties.Label != nil {
		res.StartLabel = &e.FromProperties.Label.Text
	}
	if e.Label != nil {
		res.MiddleLabel = &e.Label.Text
	}
	if e.ToProperties.Label != nil {
		res.EndLabel = &e.ToProperties.Label.Text
	}

	return res
}
