package parse

import (
	"fmt"

	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types"
)

func FindQueries(c *context.Ctx) {
	c.Queries.ParseUTML = context.DefineQuery(c, "Parse UTML", parseUTML)
	c.Queries.ParseUTMLToInternal = context.DefineQuery(c, "Parse UTML to internal repr.", convertUTMLToParseRes)
	c.Queries.Parse = context.DefineQuery(c, "Parse", parseFile)
}

// general method that switches based on context
func parseFile(c *context.Ctx, cmd types.ParseCmd) *types.ParseResult {
	switch cmd.DiagramFormat {

	case types.DiagramFormatUTML:
		utmlRes := c.Queries.ParseUTML.Get("Parse UTML", types.ParseUTMLCmd{Filepath: cmd.Filepath})
		return convertUTMLToParseRes(c, utmlRes)

	default:
		panic(fmt.Sprintf("Format '%v' not implemented for parsing!", cmd.DiagramFormat))
	}
}
