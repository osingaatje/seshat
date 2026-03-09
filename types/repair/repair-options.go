package repair

type RepairOptions struct {
	SwapEdgeLabels bool // whether to correct misplaced labels (when a student drags labels to the opposing classes)
}

func DefaultRepairOptions() RepairOptions {
	return RepairOptions{
		SwapEdgeLabels: true,
	}
}
