package repair

// for RepairOpts.SwapEdgeLabels
const DISTANCE_PERC_THRESHOLD float64 = 0.35
const DIST_INF float64 = 999999999999

// for RepairOpts.ConnectEdgEnds
const LINE_CLOSENESS_PERCENT uint8 = 8   // NOTE: From ANY point in the edge - 0-255%, 100% is the full length of the edge
const VERTEX_CLOSENESS_PERCENT uint8 = 6 // NOTE: From the BORDER of the vertex - 0-255%, 100% = length of edge of the vertex

type RepairOptions struct {
	SwapEdgeLabels  bool // whether to correct misplaced labels (when a student drags labels to the opposing classes)
	ConnectEdgeEnds bool // whether to try to connect floating edges
	FailOnError     bool // true == return nil on failure. False == return (partially) repaired diagram
}

func DefaultRepairOptions() RepairOptions {
	return RepairOptions{
		SwapEdgeLabels:  true,
		ConnectEdgeEnds: true,
		FailOnError:     false,
	}
}
