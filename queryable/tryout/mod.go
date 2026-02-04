package tryout

import (
	"fmt"
	"seshat/context"
)

func FindQueries(c *context.Ctx) {
	c.Queries.Test = context.DefineQuery(c, "Test", Test)
}

func Test(c *context.Ctx, name context.MultiKey2[string, string]) string {
	greeting := fmt.Sprintf("Hello %s %s!", name.K1, name.K2)
	return greeting
}
