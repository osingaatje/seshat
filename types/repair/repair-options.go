package repair

// for RepairOpts.SwapEdgeLabels
const DISTANCE_PERC_THRESHOLD float64 = 0.35
const DIST_INF float64 = 999999999999

// for RepairOpts.ConnectEdgEnds
const LINE_CLOSENESS_PERCENT uint8 = 10    // NOTE: From ANY point in the edge - 0-255%, 100% is the full length of the edge
const VERTEX_CLOSENESS_PERCENT uint8 = 100 // NOTE: From the CENTER of the vertex - 0-255%, 100% = size of node (avg. width and height)

type RepairOptions struct {
	SwapEdgeLabels  bool // whether to correct misplaced labels (when a student drags labels to the opposing classes)
	ConnectEdgeEnds bool // whether to try to connect floating edges
}

func DefaultRepairOptions() RepairOptions {
	return RepairOptions{
		SwapEdgeLabels:  true,
		ConnectEdgeEnds: true,
	}
}
