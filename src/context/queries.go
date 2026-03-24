package context

//---------------------------//
//   Central Queries struct  //
// - add your queries here - //
//---------------------------//

import (
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/dot"
	. "github.com/osingaatje/seshat/types/graph/internal-rep"
	. "github.com/osingaatje/seshat/types/graph/parse-result"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

type Queries struct {
	ctx *Ctx

	// queries are placed here. Note that you need to add this query to DefineQueries() and then let the Driver call the function in order to use it!
	Parse               *Query[ParseCmd, *ParseResult]
	ParseUTML           *Query[string /* file path */, *ParseResultUTML]
	ParseUTMLToParseRes *Query[*ParseResultUTML, *ParseResult]

	DisplayDiagramAsDot *Query[*ParseResult, *DotGraph]
	RepairDiagram       *Query[RepairCmd /* parse result + config */, *ParseResult]

	ConvertGraphToInternal *Query[*ParseResult, *InternalGraph]

	GradeDiagram                     *Query[GradeCmd, *GradeResult]
	SyntacticMatch                   *Query[MatchStringCmd, int]
	SemanticMatchWordnet             *Query[MatchStringCmd, MatchStringRes]
	SemanticMatchSentenceTransformer *Query[MatchStringCmd, MatchStringRes]

	Test *Query[NameCmd, string]
}
