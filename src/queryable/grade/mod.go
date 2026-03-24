package grade

import (
	"github.com/osingaatje/seshat/src/context"
)

func FindQueries(c *context.Ctx) {
	c.Queries.GradeDiagram = context.DefineQuery(c,
		"Grade Internal Diagram", gradeDiag)

	c.Queries.SyntacticMatch = context.DefineQuery(c,
		"Levenshtein distance", syntacticDist)

	c.Queries.SemanticMatchWordnet = context.DefineQuery(c,
		"WordNet Similarity Score between words", semanticMatchWordnet)

	c.Queries.SemanticMatchSentenceTransformer = context.DefineQuery(c,
		"Sentence Transformer (MiniLM L6 v2) semantic equivalence", semanticSimilarityMiniLM)
}
