package context

//---------------------------//
//   Central Queries struct  //
// - add your queries here - //
//---------------------------//

import (
	. "github.com/osingaatje/seshat/types/command"
	. "github.com/osingaatje/seshat/types/grade"
	. "github.com/osingaatje/seshat/types/graph/dot"
	. "github.com/osingaatje/seshat/types/graph/intern"
	. "github.com/osingaatje/seshat/types/graph/utml"
)

type Queries struct {
	ctx *Ctx

	// queries are placed here. Note that you need to add this query to DefineQueries() and then let the Driver call the function in order to use it!
	Parse               *Query[ParseCmd, *InternalGraph]
	ParseUTML           *Query[string /* file path */, *ParseResultUTML]
	ParseUTMLToParseRes *Query[*ParseResultUTML, *InternalGraph]

	DisplayDiagramAsDot *Query[*InternalGraph, *DotGraph]
	RepairDiagram       *Query[RepairCmd /* parse result + config */, RepairResult]

	GradeDiagram                     *Query[GradeCmd, *GradeCalculation]
	SyntacticMatch                   *Query[MatchStringCmd, int]
	SemanticMatchWordnet             *Query[MatchStringCmd, MatchStringRes]
	SemanticMatchSentenceTransformer *Query[MatchStringCmd, MatchStringRes]

	Test *Query[NameCmd, string]
}
