package command

import (
	. "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/repair"
)

type RepairCmd struct {
	Diagram    *InternalGraph
	RepairOpts RepairOptions
}

func NewRepairCmdDefOpt(p *InternalGraph) RepairCmd {
	return RepairCmd{Diagram: p, RepairOpts: DefaultRepairOptions()}
}
func NewRepairCmd(p *InternalGraph, o RepairOptions) RepairCmd {
	return RepairCmd{Diagram: p, RepairOpts: o}
}
