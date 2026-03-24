package grade

import (
	"fmt"

	"github.com/osingaatje/seshat/helper"
	sentencetransformer "github.com/osingaatje/seshat/helper/sentence-transformer"
	"github.com/osingaatje/seshat/src/context"
)

func SemanticSimilarityMiniLM(c *context.Ctx, s1 string, s2 string) (float64, error) {
	r, err := sentencetransformer.CompareSentences(c, []string{s1, s2})
	if err != nil {
		return -1, err
	}
	jsonres, _ := helper.MarshalJSON(r)
	fmt.Println(string(jsonres))
	return 0, fmt.Errorf("Not implemented")
}
