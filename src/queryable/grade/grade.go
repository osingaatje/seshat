package grade

import (
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	"github.com/osingaatje/seshat/types/grade"
)

func gradeDiag(c *context.Ctx, cmd command.GradeCmd) *grade.GradeResult {
	possibleVertexMappings := getAlternativeSolutions(c, cmd)

	if len(possibleVertexMappings) == 0 {
		c.LogErr("BRUH")
	}

	c.LogErr("TODO MAKE GRADING WORK")
	return nil
}
