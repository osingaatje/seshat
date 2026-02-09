package context

import (
	"github.com/osingaatje/seshat/context/data"
)

type Queries struct {
	ctx *Ctx

	// queries are placed here. Note that you need to add this query to DefineQueries() and then let the Driver call the function in order to use it!
	Test *Query[data.NameCmd, string]
}
