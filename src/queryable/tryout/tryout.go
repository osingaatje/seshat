package tryout

//-------------------------------------//
// An example module that uses queries //
//-------------------------------------//

import (
	"fmt"

	"github.com/osingaatje/seshat/src/context"
	. "github.com/osingaatje/seshat/types/command"
)

func FindQueries(c *context.Ctx) {
	c.Queries.Test = context.DefineQuery(c, "Test", Test)
}

func Test(c *context.Ctx, name NameCmd) string {
	greeting := fmt.Sprintf("Hello %s %s!", name.FName, name.LName)
	return greeting
}
