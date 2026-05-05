package repair

import (
	"fmt"
	"strings"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	cmd "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/graph/intern"
	"github.com/osingaatje/seshat/types/graph/shared"
)

func FindQueries(c *context.Ctx) {
	c.Queries.RepairDiagram = context.DefineQuery(
		c,
		"Repair internal representation",
		performRepairs,
	)
}

func performRepairs(c *context.Ctx, conf cmd.RepairCmd) *InternalGraph {
	c.LogPrefixAdd("Repair '%s'", conf.Diagram.Metadata.Filename)
	defer c.LogPrefixRm("Repair '%s'", conf.Diagram.Metadata.Filename)

	if conf.Diagram == nil {
		return nil // identity: repair(nil) = nil
	}

	res := conf.Diagram.Copy() // COPY THE DIAGRAM to avoid changing the original, which would break immutability

	failedEdgeCorrections := map[shared.EdgeIdentifier][]error{}
	for eId, e := range res.Edges {
		failedEdgeCorrections[eId] = []error{}

		// swap edge labels if we detect that a student has dragged them to other spots
		if conf.RepairOpts.SwapEdgeLabels {
			swapEdgeLabelsForEdge(c, e)
		}

		// Connect loose edge ends if option enabled:
		if conf.RepairOpts.ConnectEdgeEnds {
			err := tryConnectEdgeEnds(res, e)
			if err != nil {
				failedEdgeCorrections[eId] = append(failedEdgeCorrections[eId], err)
			}

			// SANITY CHECKS:
			err = verifyEdgesLinkToVertices(res, e)
			if err != nil {
				failedEdgeCorrections[eId] = append(failedEdgeCorrections[eId], err)
			}
		}
	}

	fails := helper.FilterMap(failedEdgeCorrections, func(_ shared.EdgeIdentifier, errs []error) bool {
		return len(errs) > 0
	})
	if len(fails) > 0 {
		keys := helper.Keys(fails)
		keyStrings := helper.Map(keys, func(k shared.EdgeIdentifier) string { return fmt.Sprintf("%d", k) })

		var errMsg strings.Builder
		fmt.Fprintf(&errMsg, "Repairs failed for edges [%s]", strings.Join(keyStrings, ","))
		for k, v := range fails {
			errStrings := helper.Map(v, func(e error) string { return e.Error() })
			fmt.Fprintf(&errMsg, "\nID '%d': [%s]", k, strings.Join(errStrings, ","))
		}
		c.LogErr(errMsg.String())
		return nil
	}

	return res
}

func verifyEdgesLinkToVertices(g *InternalGraph, e *InternalEdge) error {
	var hasFromVertex bool = false
	var hasToVertex bool = false
	var hasFromEdge bool = false
	var hasToEdge bool = false

	if e.FromId != nil {
		_, hasFromVertex = g.Vertices[*e.FromId]
	}
	if e.ToId != nil {
		_, hasToVertex = g.Vertices[*e.ToId]
	}
	if e.FromEdgeId != nil {
		_, hasFromEdge = g.Edges[*e.FromEdgeId]
	}
	if e.ToEdgeId != nil {
		_, hasToEdge = g.Edges[*e.ToEdgeId]
	}

	// if the edge has both (or neither) a from vertex/edge or a to vertex/edge, then we did something wrong
	if hasFromVertex == hasFromEdge || hasToVertex == hasToEdge {
		return fmt.Errorf("Edge '%d' has either no or both a starting or ending vertex/edge", e.Id)
	}
	return nil
}
