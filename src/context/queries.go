package context

//---------------------------//
//   Central Queries struct  //
// - add your queries here - //
//---------------------------//

import (
	"github.com/osingaatje/seshat/types"
)

type Queries struct {
	ctx *Ctx

	// queries are placed here. Note that you need to add this query to DefineQueries() and then let the Driver call the function in order to use it!
	Parse               *Query[data.ParseCmd, *data.ParseResult]
	ParseUTML           *Query[data.ParseUTMLCmd, *data.ParseResultUTML]
	ParseUTMLToInternal *Query[*data.ParseResultUTML, *data.ParseResult]

	Test *Query[data.NameCmd, string]
}
