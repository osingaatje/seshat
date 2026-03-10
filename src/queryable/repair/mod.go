package repair

import (
	"github.com/osingaatje/seshat/src/context"
	cmd "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/parse-result"
)

func FindQueries(c *context.Ctx) {
	c.Queries.RepairDiagram = context.DefineQuery(
		c,
		"Repair internal representation",
		performRepairs)
}

func performRepairs(c *context.Ctx, conf cmd.RepairCmd) *ParseResult {
	if conf.Diagram == nil {
		return nil // repair(nil) = nil
	}

	res := conf.Diagram.Copy() // COPY THE DIAGRAM to avoid changing the original, which would break immutability

	// swap edge labels if we detect that a student has dragged them to other spots
	if conf.RepairOpts.SwapEdgeLabels {
		swapEdgeLabels(c, res)
	}

	// If more options become available, add them here!

	return res
}
