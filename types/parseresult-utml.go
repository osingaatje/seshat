package types

import (
	"github.com/osingaatje/seshat/helper"
)

type UTMLClassType string

type UTMLVisibility string

const (
	UTMLClassTypeClass     UTMLClassType = "class"
	UTMLClassTypeInterface UTMLClassType = "interface"
	UTMLClassTypeAbstract  UTMLClassType = "abstract"

	UTMLVisibilityPrivate   UTMLVisibility = "private"
	UTMLVisibilityPublic    UTMLVisibility = "public"
	UTMLVisibilityProtected UTMLVisibility = "protected"
	UTMLVisibilityPackage   UTMLVisibility = "package"
)

func (p *ParseResultUTML) String() string {
	res, err := helper.MarshalJSON(p)
	if err != nil {
		panic("somehow I can't marshal UTML parse results to JSON? FIX!")
	}
	return string(res)
}

type ParseResultUTML struct {
	Edges []ParseResultUTMLEdge `json:"edges"`

	Nodes []ParseResultUTMLNode `json:"nodes"`
}

type ParseResultUTMLEdge struct {
	StartPosition int `json:"startPosition"`
	EndPosition   int `json:"endPosition"`

	StartLabel      *UTMLEdgeLabel      `json:"startLabel,omitempty"`  // text
	MiddleLabel     *UTMLEdgeLabel      `json:"middleLabel,omitempty"` // text
	EndLabel        *UTMLEdgeLabel      `json:"endLabel,omitempty"`    // text
	StartStyle      *UTMLArrowHeadStyle `json:"startStyle,omitempty"`  // arrow head style?
	EndStyle        *UTMLArrowHeadStyle `json:"endStyle,omitempty"`    // arrow head style?
	LineStyle       *UTMLLineStyle      `json:"lineStyle,omitempty"`   // line styling
	LineType        *int                `json:"lineType,omitempty"`    // line styling
	MiddlePositions []UTMLXY            `json:"middlePositions"`       // no clue what this is
	StartNodeId     int                 `json:"startNodeId"`           // node pointer
	EndNodeId       int                 `json:"endNodeId"`             // node pointer
}

type UTMLLineStyle int

const (
	UTMLLineStyleFilled UTMLLineStyle = 0
	UTMLLineStyleDotted               = iota
	UTMLLineStyleDashed
)

type UTMLArrowHeadStyle int

const (
	UTMLArrowStyleNone UTMLArrowHeadStyle = iota
	UTMLArrowStyleSmallFilledArrow
	UTMLArrowStyleFilledDiamond
	UTMLArrowStyleUnfilledDiamond
	UTMLArrowStyleLargeUnfilledArrow
)

type ParseResultUTMLNode struct {
	Type            *UTMLClassType `json:"type,omitempty"`
	Width           *int           `json:"width,omitempty"`
	Height          *int           `json:"height,omitempty"`
	Position        *UTMLXY        `json:"position,omitempty"`
	Text            *string        `json:"text,omitempty"`
	HasDoubleBorder *bool          `json:"hasDoubleBorder,omitempty"`
	StyleObject     *struct {
		Fill          string  `json:"fill"`
		Stroke        string  `json:"stroke"`
		StrokeWidth   float32 `json:"stroke-width"`
		FillOpacity   float32 `json:"fill-opacity"`
		StrokeOpacity float32 `json:"stroke-opacity"`
	} `json:"styleObject,omitempty"`
	ClassType  *string             `json:"classType,omitempty"`
	Attributes []UTMLFieldOrMethod `json:"attributes,omitempty"`
	Methods    []UTMLFieldOrMethod `json:"methods,omitempty"`
}

type UTMLEdgeLabel struct {
	Offset UTMLXY `json:"offset"`

	EdgeLocation int    `json:"edgeLocation"`
	Value        string `json:"value"`
}

type UTMLFieldOrMethod struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Visibility UTMLVisibility `json:"visibility"`
}

type UTMLXY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

/* EXAMPLE
{
  "edges": [
    {
      "startPosition": 6,
      "endPosition": 2,
      "startLabel": {
        "offset": {
          "x": -19.970314158618923,
          "y": -1.0892898631973957
        },
        "edgeLocation": 0,
        "value": "1"
      },
      "middleLabel": {
        "offset": {
          "x": 0,
          "y": 0
        },
        "edgeLocation": 1,
        "value": "< commisions"
      },
      "endLabel": {
        "offset": {
          "x": 19.970314158618923,
          "y": 1.0892898631973957
        },
        "edgeLocation": 2,
        "value": "1..*"
      },
      "startStyle": 0,
      "endStyle": 0,
      "lineStyle": 0,
      "lineType": 5,
      "middlePositions": [],
      "startNodeId": 1,
      "endNodeId": 0
    },
		...
	],
	"nodes": [
    {
      "type": "ClassNode",
      "width": 190,
      "height": 128,
      "position": {
        "x": 1010,
        "y": 240
      },
      "text": "Project",
      "hasDoubleBorder": false,
      "styleObject": {
        "fill": "white",
        "stroke": "black",
        "stroke-width": 2,
        "fill-opacity": 1,
        "stroke-opacity": 0.75
      },
      "classType": "class",
      "attributes": [
        {
          "name": "startDate",
          "type": "String",
          "visibility": "private"
        },
        {
          "name": "deadline",
          "type": "String",
          "visibility": "private"
        },
        {
          "name": "budget",
          "type": "double",
          "visibility": "private"
        },
        {
          "name": "id",
          "type": "int",
          "visibility": "private"
        }
      ],
      "methods": []
    },
		...
	]
}
*/
