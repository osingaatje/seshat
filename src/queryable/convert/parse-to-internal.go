package convert

import (
	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/graph/internal-rep"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
)

func convertParseResToInternal(c *context.Ctx, p *ParseResult) *InternalGraph {
	res, err := p.ToInternal()
	if err != nil {
		c.LogErr("Error while parsing parse result to internal representation: %s", err.Error())
		return nil
	}

	return res
}
