package utml

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/osingaatje/seshat/helper"
	. "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/graph/shared"
)

var SKIPPED_VERTEX_TYPES []string = []string{"CommentNode"}

type UTMLVisibility string

const (
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
	// metadata
	Metadata GraphMetadata `json:"-"` // we don't serialise it because then it's not the original JSON file.

	Edges []ParseResultUTMLEdge `json:"edges"`

	Nodes []ParseResultUTMLNode `json:"nodes"`
}

type ParseResultUTMLEdge struct {
	StartPosition UTMLEdgeXYOrOffsetPosition `json:"startPosition"` // indicates either an offset to the node (StartNodeId) or an absolute position
	EndPosition   UTMLEdgeXYOrOffsetPosition `json:"endPosition"`   // indicates either an offset to the node (EndNodeId) or absolute pos.

	StartLabel      *UTMLEdgeLabel     `json:"startLabel,omitempty"`  // text
	MiddleLabel     *UTMLEdgeLabel     `json:"middleLabel,omitempty"` // text
	EndLabel        *UTMLEdgeLabel     `json:"endLabel,omitempty"`    // text
	StartStyle      UTMLArrowHeadStyle `json:"startStyle"`            // arrow head style?
	EndStyle        UTMLArrowHeadStyle `json:"endStyle"`              // arrow head style?
	LineStyle       UTMLLineStyle      `json:"lineStyle"`             // line styling
	LineType        int                `json:"lineType"`              // line styling
	MiddlePositions []UTMLXY           `json:"middlePositions"`       // no clue what this is
	StartNodeId     *int16             `json:"startNodeId,omitempty"` // node pointer
	EndNodeId       *int16             `json:"endNodeId,omitempty"`   // node pointer
}

func (p ParseResultUTMLEdge) Hash() int { // at least 32 bits in size so shifting 16-bit vals should be fine.
	jsonbytes, err := helper.MarshalJSON(p)
	if err != nil {
		panic("AAA why can I not convert to JSON? stupid")
	}
	h := fnv.New32a()
	h.Write(jsonbytes)
	return int(h.Sum32())
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
	Type            string `json:"type"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Position        UTMLXY `json:"position"`
	Text            string `json:"text"`
	HasDoubleBorder bool   `json:"hasDoubleBorder"`
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

type UTMLEdgeXYOrOffsetPosition struct {
	Value any // UTMLEdgeOffsetPosition | UTMLXY
}

func (pos UTMLEdgeXYOrOffsetPosition) MarshalJSON() ([]byte, error) {
	switch val := pos.Value.(type) {
	case UTMLEdgeOffsetPosition:
		return json.Marshal(int(val))
	case UTMLXY:
		return json.Marshal(val)
	}
	return nil, fmt.Errorf("Cannot marshal UTMLEdgeXYOrOffset if it's not an offset or XY position!")
}

func (pos *UTMLEdgeXYOrOffsetPosition) UnmarshalJSON(data []byte) error {
	// first try UTMLEdgeOffsetPosition
	var offsetVal UTMLEdgeOffsetPosition
	if err := json.Unmarshal(data, &offsetVal); err == nil {
		pos.Value = offsetVal
		return nil
	}

	// otherwise try X/Y val
	var xyVal UTMLXY
	if err := json.Unmarshal(data, &xyVal); err == nil {
		pos.Value = xyVal
		return nil
	}

	panic("Edge offset / position was neither UTMLEdgeOffsetPosition nor UTMLXY value - bug in the code")
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
type UTMLEdgeOffsetPosition int

const (
	EdgePosTopCenter    UTMLEdgeOffsetPosition = 0
	EdgePosTopRight     UTMLEdgeOffsetPosition = 1
	EdgePosMiddleRight  UTMLEdgeOffsetPosition = 2
	EdgePosBottomRight  UTMLEdgeOffsetPosition = 3
	EdgePosBottomCenter UTMLEdgeOffsetPosition = 4
	EdgePosBottomLeft   UTMLEdgeOffsetPosition = 5
	EdgePosMiddleLeft   UTMLEdgeOffsetPosition = 6
	EdgePosTopLeft      UTMLEdgeOffsetPosition = 7
)
