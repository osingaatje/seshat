package tryout

import (
	"fmt"
	"seshat/context"
)

func FindQueries(c *context.Ctx) {
	c.Queries.Test = context.DefineQuery(c, "Test", Test)
}

func Test(c *context.Ctx, name string) string {
	greeting := fmt.Sprintf("Hello %s!", name)
	return greeting
}
