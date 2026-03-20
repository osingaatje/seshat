package convert

import (
	"fmt"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
)

func FindQueries(c *context.Ctx) {
	c.Queries.Parse = context.DefineQuery(c,
		"Parse file", parseFile)
	c.Queries.ParseUTML = context.DefineQuery(c,
		"Parse UTML", parseUTML)

	c.Queries.ParseUTMLToParseRes = context.DefineQuery(c,
		"Parse UTML to parse result", convertUTMLToParseRes)

	c.Queries.ConvertGraphToInternal = context.DefineQuery(c,
		"Parse Result -> Internal graph", convertParseResToInternal)
}

// general method that switches based on context
func parseFile(c *context.Ctx, cmd ParseCmd) *ParseResult {
	switch cmd.DiagramFormat {

	case DiagramFormatUTML:
		utmlRes := c.Queries.ParseUTML.Get("Parse UTML", cmd.Filepath)
		return convertUTMLToParseRes(c, utmlRes)

	default:
		panic(fmt.Sprintf("Format '%v' not implemented for parsing!", cmd.DiagramFormat))
	}
}
