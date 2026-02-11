package parse

import (
	"fmt"

	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/context/data"
)

func FindQueries(c *context.Ctx) {
	c.Queries.ParseUTML = context.DefineQuery(c, "Parse UTML", parseUTML)
	c.Queries.Parse = context.DefineQuery(c, "Parse", parseFile)
}

// general method that switches based on context
func parseFile(c *context.Ctx, cmd data.ParseCmd) *data.ParseResult {
	switch cmd.DiagramFormat {

	case data.DiagramFormatUTML:
		utmlRes := c.Queries.ParseUTML.Get("Parse UTML", data.ParseUTMLCmd{Filepath: cmd.Filepath})
		return parseToParseResult(utmlRes)

	default:
		panic(fmt.Sprintf("Format '%s' not implemented for parsing!", cmd.DiagramFormat))
	}
}
