package repair

// for RepairOpts.SwapEdgeLabels
const DISTANCE_PERC_THRESHOLD float64 = 0.35
const DIST_INF float64 = 999999999999

// for RepairOpts.ConnectEdgEnds
const LINE_CLOSENESS_PERCENT uint8 = 8   // NOTE: From ANY point in the edge - 0-255%, 100% is the full length of the edge
const VERTEX_CLOSENESS_PERCENT uint8 = 6 // NOTE: From the BORDER of the vertex - 0-255%, 100% = length of edge of the vertex

type RepairOptions struct {
	/*
	 * whether to correct misplaced labels (when a student drags labels to the opposing classes)
	 */
	SwapEdgeLabels bool

	// whether to try to connect 'floating' edges:  [V1]  -------[V2]   becomes   [V1]------[V2]
	ConnectEdgeEnds bool

	/* whether to simplify edge-to-edge-to-vertex connections IF THEY HAVE A CLEAR DIRECTION:
		  		[V1]<--*----*---[V2]
						\
						 ----[V3]
				would become:
				[V1]<----------[V2]
	              ^
				   \________[V3]
	*/
	SimplifyDirectedEdges bool

	// if any reparations fail, whether to return a 'nil' result and log errors etc. (TRUE), or whether to only report errors and not scream (FALSE).
	FailOnError bool
}

func DefaultRepairOptions() RepairOptions {
	return RepairOptions{
		SwapEdgeLabels:        true,
		ConnectEdgeEnds:       true,
		SimplifyDirectedEdges: true,
		FailOnError:           false,
	}
}
