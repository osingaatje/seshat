package utml

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
	StartPosition UTMLEdgeEndPosition `json:"startPosition"` // indicates to which end of the node the edge connects
	EndPosition   UTMLEdgeEndPosition `json:"endPosition"`   // indicates to which end of the node the edge connects

	StartLabel      *UTMLEdgeLabel     `json:"startLabel,omitempty"`  // text
	MiddleLabel     *UTMLEdgeLabel     `json:"middleLabel,omitempty"` // text
	EndLabel        *UTMLEdgeLabel     `json:"endLabel,omitempty"`    // text
	StartStyle      UTMLArrowHeadStyle `json:"startStyle"`            // arrow head style?
	EndStyle        UTMLArrowHeadStyle `json:"endStyle"`              // arrow head style?
	LineStyle       UTMLLineStyle      `json:"lineStyle"`             // line styling
	LineType        int                `json:"lineType"`              // line styling
	MiddlePositions []UTMLXY           `json:"middlePositions"`       // no clue what this is
	StartNodeId     int                `json:"startNodeId"`           // node pointer
	EndNodeId       int                `json:"endNodeId"`             // node pointer
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
	Type            UTMLClassType `json:"type"`
	Width           int           `json:"width"`
	Height          int           `json:"height"`
	Position        UTMLXY        `json:"position"`
	Text            string        `json:"text"`
	HasDoubleBorder bool          `json:"hasDoubleBorder"`
	StyleObject     *struct {
		Fill          string  `json:"fill"`
		Stroke        string  `json:"stroke"`
		StrokeWidth   uint16  `json:"stroke-width"`
		FillOpacity   float32 `json:"fill-opacity"`
		StrokeOpacity float32 `json:"stroke-opacity"`
	} `json:"styleObject,omitempty"`
	ClassType  *string             `json:"classType,omitempty"`
	Attributes []UTMLFieldOrMethod `json:"attributes"`
	Methods    []UTMLFieldOrMethod `json:"methods"`
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
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

/*
 *     7     0     1
 *     |----------|
 *     |          |
 *   6 |          | 2
 *     |          |
 *     |----------|
 *    5     4      3
 */
type UTMLEdgeEndPosition int

const (
	EdgePosTopCenter    UTMLEdgeEndPosition = 0
	EdgePosTopRight     UTMLEdgeEndPosition = 1
	EdgePosMiddleRight  UTMLEdgeEndPosition = 2
	EdgePosBottomRight  UTMLEdgeEndPosition = 3
	EdgePosBottomCenter UTMLEdgeEndPosition = 4
	EdgePosBottomLeft   UTMLEdgeEndPosition = 5
	EdgePosMiddleLeft   UTMLEdgeEndPosition = 6
	EdgePosTopLeft      UTMLEdgeEndPosition = 7
)
