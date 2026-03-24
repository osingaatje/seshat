package sentencetransformer

import (
	"math"
	"os"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
	"github.com/osingaatje/seshat/src/context"
)

// TODO: Perhaps use a combination of MiniLM-L6-v2 and msmarco-MiniLM (see Ramachandran, 2025)

const DEFAULT_SENTENCE_TRANSFORMER_MODEL_PATH = "./al-MiniLM-L6-v2"

// loads env SENTENCE_TRANSFORMER_MODEL_PATH or falls back to default
func modelPath(c *context.Ctx) string {
	path := os.Getenv("SENTENCE_TRANSFORMER_MODEL_PATH")
	if path == "" {
		c.LogWarn("No .env with SENTENCE_TRANSFORMER_MODEL_PATH specified, falling back to default: %s", DEFAULT_SENTENCE_TRANSFORMER_MODEL_PATH)
		path = DEFAULT_SENTENCE_TRANSFORMER_MODEL_PATH
	}
	return path
}

var hugotMiniLML6v2Pipeline *pipelines.FeatureExtractionPipeline = nil
var hugotMiniLML6v2PipelineErr error = nil

func CompareSentences(c *context.Ctx, s []string) ([][]float32, error) {
	if hugotMiniLML6v2PipelineErr != nil {
		return nil, c.LogErrAndReturn("Error while initialising Hugot Pipeline: %s", hugotMiniLML6v2PipelineErr.Error())
	}
	if hugotMiniLML6v2Pipeline == nil {
		// init session and pipline etc.
		session, hugotMiniLML6v2PipelineErr := hugot.NewGoSession()
		if hugotMiniLML6v2PipelineErr != nil {
			return nil, c.LogErrAndReturn("Could not init Hugot semantic analysis session: %s", hugotMiniLML6v2PipelineErr.Error())
		}

		// don't destroy the session until the end of the program
		//defer func(c *context.Ctx, session *hugot.Session) {
		//	err := session.Destroy()
		//	if err != nil {
		//		c.LogErr("Could not destroy Hugot session! Err=%s", err.Error())
		//	}
		//}(c, session)

		dopts := hugot.NewDownloadOptions()
		dopts.OnnxFilePath = "onnx/model.onnx"

		modelPath, hugotMiniLML6v2PipelineErr := hugot.DownloadModel("sentence-transformers/all-MiniLM-L6-v2", modelPath(c), dopts)
		if hugotMiniLML6v2PipelineErr != nil {
			return nil, c.LogErrAndReturn("Could not download MiniLM-L6-v2: %s", hugotMiniLML6v2PipelineErr.Error())
		}

		config := hugot.FeatureExtractionConfig{
			ModelPath:    modelPath,
			Name:         "SemanticSimilarity-MiniLM-L6-v2",
			OnnxFilename: "model.onnx",
		}
		hugotMiniLML6v2Pipeline, hugotMiniLML6v2PipelineErr = hugot.NewPipeline(session, config)
		if hugotMiniLML6v2PipelineErr != nil {
			return nil, c.LogErrAndReturn("Could not init MiniLM-L6-v2 similarity pipeline! Err: %s", hugotMiniLML6v2PipelineErr.Error())
		}
	}

	result, hugotMiniLML6v2Pipeline := hugotMiniLML6v2Pipeline.RunPipeline(s)
	if hugotMiniLML6v2Pipeline != nil {
		return nil, c.LogErrAndReturn("Could not run pipeline: %s", hugotMiniLML6v2Pipeline.Error())
	}
	return result.Embeddings, nil
}

func CosineSimilarity(v1 []float32, v2 []float32) (cosine float64) {
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
