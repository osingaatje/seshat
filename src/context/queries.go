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
	Parse               *Query[types.ParseCmd, *types.ParseResult]
	ParseUTML           *Query[types.ParseUTMLCmd, *types.ParseResultUTML]
	ParseUTMLToInternal *Query[*types.ParseResultUTML, *types.ParseResult]

	Test *Query[types.NameCmd, string]
}
