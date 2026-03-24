package sentencetransformer

import (
	"os"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
	"github.com/osingaatje/seshat/src/context"
)

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

func CompareSentences(c *context.Ctx, s []string) (*pipelines.FeatureExtractionOutput, error) {
	session, err := hugot.NewGoSession()
	if err != nil {
		return nil, err
	}

	defer func(c *context.Ctx, session *hugot.Session) {
		err := session.Destroy()
		if err != nil {
			c.LogErr("Could not destroy Hugot session! Err=%s", err.Error())
		}
	}(c, session)

	dopts := hugot.NewDownloadOptions()
	dopts.OnnxFilePath = "onnx/model.onnx"

	modelPath, err := hugot.DownloadModel("sentence-transformers/all-MiniLM-L6-v2", modelPath(c), dopts)
	if err != nil {
		return nil, c.LogErrAndReturn("Could not download MiniLM-L6-v2: %s", err.Error())
	}

	config := hugot.FeatureExtractionConfig{
		ModelPath:    modelPath,
		Name:         "SemanticSimilarity-MiniLM-L6-v2",
		OnnxFilename: "model.onnx",
	}
	pipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		return nil, c.LogErrAndReturn("Could not init semantic similarity pipeline! Err: %s", err.Error())
	}

	return pipeline.RunPipeline(s)
}
