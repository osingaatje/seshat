package displaygraph

import (
	"fmt"
	"strings"

	. "github.com/osingaatje/seshat/types/generic"
)

// example:
// https://dreampuf.github.io/GraphvizOnline/?engine=fdp#%2F%2F%20test%0Agraph%20test%20%7B%0A%20%20%20%20%2F%2F%20Graph%20attributes%0A%20%20%20%20overlap%3Dfalse%3B%0A%20%20%20%20splines%3Dtrue%3B%0A%20%20%20%20node%20%5Bshape%3Dsquare%2C%20style%3Dfilled%2C%20fillcolor%3Dlightblue%2C%20fontname%3D%22Arial%22%5D%3B%0A%20%20%20%20edge%20%5Bcolor%3Dgray50%2C%20penwidth%3D1.5%5D%3B%0A%0A%20%20%20%20%2F%2F%20Nodes%20(people)%0A%20%20%20%20Alice%5Bpos%3D%2210%2C120%22%5D%3B%0A%20%20%20%20Bob%5Bpos%3D%2220%2C100%22%5D%3B%0A%20%20%20%20Carol%5Bpos%3D%2230%2C100%22%5D%3B%0A%20%20%20%20Dave%5Bpos%3D%2240%2C100%22%5D%3B%0A%20%20%20%20Eve%5Bpos%3D%2280%2C100%22%5D%3B%0A%0A%20%20%20%20%2F%2F%20Edges%20(friendships)%0A%20%20%20%20Alice%20--%20Bob%20%5Blabel%3D%22test%22%2Ctaillabel%3D%22testtail%22%2Cheadlabel%3D%22testhead%22%5D%3B%0A%20%20%20%20Alice%20--%20Carol%3B%0A%20%20%20%20Bob%20--%20Carol%3B%0A%20%20%20%20Bob%20--%20Dave%3B%0A%20%20%20%20Carol%20--%20Dave%3B%0A%20%20%20%20Dave%20--%20Eve%3B%0A%20%20%20%20Eve%20--%20Alice%3B%0A%7D
type DotGraph struct {
	Name          string
	GraphSettings DotGraphSettings
	NodeSettings  DotNodeSettings
	EdgeSettings  DotEdgeSettings

	Nodes []DotNode
	Edges []DotEdge
}

func (d DotGraph) New(name string) DotGraph {
	return DotGraph{
		Name:          name,
		GraphSettings: DotGraphSettings{},
		NodeSettings:  DotNodeSettings{},
		EdgeSettings:  DotEdgeSettings{},

		Nodes: []DotNode{},
		Edges: []DotEdge{},
	}
}

func (g *DotGraph) String() string {
	res := "digraph " + g.Name + " {\n"

	res += g.GraphSettings.String()
	res += "\t" + g.NodeSettings.String()
	res += "\t" + g.EdgeSettings.String()

	for _, n := range g.Nodes {
		res += "\n\t" + n.String() + ";"
	}

	for _, e := range g.Edges {
		res += "\n\t" + e.String()
	}

	res += "\n}"

	return res
}

type DotNode struct {
	Text     string
	NodeOpts DotNodeOptions
}

func (d *DotNode) ToStringNoOpts() string {
	return "\"" + d.Text + "\""
}

func (d *DotNode) String() string {
	return d.ToStringNoOpts() + d.NodeOpts.String()
}

type DotNodeOptions struct {
	Pos       *Vector2D
	NoJustify bool
	Style     string
}

func (o *DotNodeOptions) String() string {
	res := "["
	elems := []string{}
	// per case:
	if o.Pos != nil {
		elems = append(elems, fmt.Sprintf("pos=\"%.1f,%.1f!\"", o.Pos.X, o.Pos.Y))
	}
	if o.NoJustify {
		elems = append(elems, "nojustify=true")
	}
	if o.Style != "" {
		elems = append(elems, "style=\""+o.Style+"\"")
	}

	res += strings.Join(elems, ",")
	res += "]"
	return res
}

type DotEdge struct {
	IdText      string
	FromText    string
	ToText      string
	MiddleLabel string  // maps to "label"
	StartLabel  *string // maps to "taillabel"
	EndLabel    *string // maps to "headlabel"

	EdgeOpts DotEdgeOptions
}

func (e *DotEdge) String() string {
	fromTxt := e.FromText
	toTxt := e.ToText
	if e.FromText == "" {
		fromTxt = "Missing start vertex or edge"
	}
	if e.ToText == "" {
		toTxt = "Missing end vertex or edge"
	}

	fromTxt = "\"" + fromTxt + "\""
	toTxt = "\"" + toTxt + "\""

	// We add a 'ghost' element (vertex) in the center to be able to connect edges to edges.
	// looks something like this:
	//  <VERTEX1> ------- (ghost vertex) --------- <VERTEX2>
	edgeText := e.IdText

	// ghost node mustn't have quotes
	ghostNode := DotNode{
		Text: edgeText,
		NodeOpts: DotNodeOptions{
			Style: "dotted",
		},
	}

	// wrap in qutes for other formatting:
	edgeText = "\"" + edgeText + "\""
	firstPartEdgeOpts := e.EdgeOpts
	secondPartEdgeOpts := e.EdgeOpts
	firstPartEdgeOpts.EndShape = ""
	secondPartEdgeOpts.StartShape = ""

	f := formatEdge(fromTxt, edgeText, e.StartLabel, "", &e.MiddleLabel, firstPartEdgeOpts)
	m := ghostNode.String()
	t := formatEdge(edgeText, toTxt, nil, "", e.EndLabel, secondPartEdgeOpts)

	return f + "; " + m + "; " + t + ";"
}

type DotEdgeOptions struct {
	StartShape string
	Style      string
	EndShape   string
}

func formatEdge(start string, end string, startlbl *string, middlelbl string, endlbl *string, opts DotEdgeOptions) string {
	res := fmt.Sprintf("%s -> %s", start, end)

	elems := []string{}
	if startlbl != nil && (*startlbl) != "" {
		elems = append(elems, fmt.Sprintf("taillabel=\"%s\"", *startlbl))
	}
	if middlelbl != "" {
		elems = append(elems, fmt.Sprintf("label=\"%s\"", middlelbl))
	}
	if endlbl != nil && (*endlbl) != "" {
		elems = append(elems, fmt.Sprintf("headlabel=\"%s\"", *endlbl))
	}

	// regular edge options
	if opts.Style != "" {
		elems = append(elems, "style=\""+opts.Style+"\"")
	}

	arrowheadEndshape := opts.EndShape
	if arrowheadEndshape == "" {
		arrowheadEndshape = "none"
	}
	elems = append(elems, fmt.Sprintf("arrowhead=\"%s\"", arrowheadEndshape))

	if opts.StartShape != "" && opts.StartShape != "none" {
		// according to graphviz docs we use 'dir=both' in combination with 'arrowtail=...' (https://graphviz.org/docs/attrs/arrowtail/)
		elems = append(elems, fmt.Sprintf("arrowtail=\"%s\"", opts.StartShape))
		elems = append(elems, "dir=\"both\"")
	}

	res += "["
	res += strings.Join(elems, ",")
	res += "]"

	return res
}

type DotGraphSettings struct {
	Overlap bool
	Splines bool
}

func (s *DotGraphSettings) String() string {
	res := fmt.Sprintf("\toverlap=%t;\n", s.Overlap)
	res += fmt.Sprintf("\tsplines=%t;\n", s.Splines)
	res += "\n"
	return res
}

type DotNodeSettings struct {
	Shape     string
	Style     string
	FillColor string
	FontName  string
}

func (s *DotNodeSettings) String() string {
	if s.Shape == "" && s.Style == "" && s.FillColor == "" && s.FontName == "" {
		return ""
	}

	res := "node ["
	elems := []string{}

	if s.Shape != "" {
		elems = append(elems, fmt.Sprintf("shape=%s", s.Shape))
	}
	if s.Style != "" {
		elems = append(elems, fmt.Sprintf("style=%s", s.Style))
	}
	if s.FillColor != "" {
		elems = append(elems, fmt.Sprintf("fillcolor=%s", s.FillColor))
	}
	if s.FontName != "" {
		elems = append(elems, fmt.Sprintf("fontname=%s", s.FontName))
	}

	res += strings.Join(elems, ",")
	res += "];\n"
	return res
}

type DotEdgeSettings struct {
	Color    string
	PenWidth *float32
	Style    string
}

func (s *DotEdgeSettings) String() string {
	if s.Color == "" && s.PenWidth == nil {
		return ""
	}

	res := "edge ["
	elems := []string{}
	if s.Color != "" {
		elems = append(elems, fmt.Sprintf("color=%s", s.Color))
	}
	if s.PenWidth != nil {
		elems = append(elems, fmt.Sprintf("penwidth=%.2f", *s.PenWidth))
	}
	if s.Style != "" {
		elems = append(elems, fmt.Sprintf("style=\"%s\"", s.Style))
	}
	res += strings.Join(elems, ",")
	res += "];\n"

	return res
}
