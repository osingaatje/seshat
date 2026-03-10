package grade

import (
	"github.com/osingaatje/seshat/src/context"
)

func FindQueries(c *context.Ctx) {
	c.Queries.GradeDiagram = context.DefineQuery(c,
		"Grade Internal Diagram", gradeDiag)
}
