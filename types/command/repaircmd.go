package command

import (
	"fmt"
	"strings"

	"github.com/osingaatje/seshat/helper"
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

type RepairResult struct {
	Diagram *InternalGraph
	Errors  []error
}

func (r RepairResult) String() string {
	if r.Diagram == nil {
		return fmt.Sprintf("{ Diagam: nil. Errors: %s }", r.Error())
	}
	return fmt.Sprintf("{ Diagram: '%s'. Errors: %s }", r.Diagram.Metadata.Filename, r.Error())
}

func (r RepairResult) Error() string {
	errStrings := helper.Map(r.Errors, func(e error) string { return e.Error() })
	return fmt.Sprintf(strings.Join(errStrings, ","))
}
