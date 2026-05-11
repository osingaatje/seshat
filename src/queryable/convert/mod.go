package convert

import (
	"fmt"

	"github.com/osingaatje/seshat/src/context"
	u "github.com/osingaatje/seshat/src/queryable/convert/utml"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/graph/intern"
)

func FindQueries(c *context.Ctx) {
	c.Queries.Parse = context.DefineQuery(c,
		"Parse file", parseFile)
	c.Queries.ParseUTML = context.DefineQuery(c,
		"Parse UTML", u.ParseUTML)

	c.Queries.ParseUTMLToParseRes = context.DefineQuery(c,
		"Parse UTML to parse result", u.ConvertUTMLToParseRes)
}

// general method that switches based on context
func parseFile(c *context.Ctx, cmd ParseCmd) *InternalGraph {
	switch cmd.DiagramFormat {

	case DiagramFormatUTML:
		utmlRes := c.Queries.ParseUTML.Get("Parse UTML", cmd.Filepath)
		return u.ConvertUTMLToParseRes(c, utmlRes)

	default:
		panic(fmt.Sprintf("Format '%v' not implemented for parsing!", cmd.DiagramFormat))
	}
}
