package command

import (
	. "github.com/osingaatje/seshat/types/parse-result"
	. "github.com/osingaatje/seshat/types/repair"
)

type RepairCmd struct {
	Diagram    *ParseResult
	RepairOpts RepairOptions
}

func NewRepairCmd(p *ParseResult, o RepairOptions) RepairCmd {
	return RepairCmd{Diagram: p, RepairOpts: o}
}
