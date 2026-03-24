package grade

import (
	"math"

	sentencetransformer "github.com/osingaatje/seshat/helper/sentence-transformer"
	"github.com/osingaatje/seshat/src/context"
)

func SemanticSimilarityMiniLM(c *context.Ctx, s1 string, s2 string) (float64, error) {
	r, err := sentencetransformer.CompareSentences(c, []string{s1, s2})
	if err != nil {
		return -1, err
	}
	if len(r) != 2 {
		return 0, c.LogErrAndReturn("Did not expect more or less than two result vectors from comparing two strings!")
	}

	sim := cosineSimilarity(r[0], r[1])
	return sim, nil
}

func cosineSimilarity(v1 []float32, v2 []float32) (cosine float64) {
	length := len(v1)
	if len(v2) > len(v1) {
		length = len(v2)
	}

	// inspired from article: https://medium.com/@nishshekh/cosine-similarity-for-embedding-vectors-0af52eef8b74
	// and Wikipedia: https://en.wikipedia.org/wiki/Cosine_similarity
	// and https://github.com/gaspiman/cosine_similarity/blob/master/cosine.go
	// and https://github.com/kabychow/go-cosinesimilarity/blob/master/cosine_similarity.go

	var dotAB float64 = 0 // A•B
	var lenA float64 = 0  // ||A||
	var lenB float64 = 0  // ||B||

	for i := 0; i < length; i++ {
		x_i := float64(v1[i])
		y_i := float64(v2[i])

		dotAB += x_i * y_i
		lenA += math.Pow(x_i, 2)
		lenB += math.Pow(y_i, 2)
	}
	if lenA == 0 || lenB == 0 {
		return 0
	}

	return dotAB / (math.Sqrt(lenA) * math.Sqrt(lenB)) // (A•B) / (||A|| ||B||)
}
