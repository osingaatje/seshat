package parse

import (
	"encoding/json"
	"os"

	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/context/data"
)

func parseUTML(c *context.Ctx, cmd data.ParseUTMLCmd) *data.ParseResultUTML {
	r, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		c.LogError("Error occurred while reading file! Err='%s'", err.Error())
		return nil
	}

	var jsonRes *data.ParseResultUTML

	err = json.Unmarshal(r, &jsonRes)
	if err != nil {
		c.LogError("Could not marshal file '%s' to a UTML Parse Result! Err=%s", cmd.Filepath, err.Error())
		return nil
	}

	return jsonRes
}

func parseToParseResult(utml *data.ParseResultUTML) *data.ParseResult {
	return &data.ParseResult{} // not implemented yet
}
