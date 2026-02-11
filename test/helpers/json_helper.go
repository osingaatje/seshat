package helpers

import (
	"bytes"
	"encoding/json"
)

func IndentJSON(fileContents []byte) (string, error) {
	compact := new(bytes.Buffer)
	// first compact it
	err := json.Compact(compact, fileContents)
	if err != nil {
		return "", err
	}

	// now indent the shit out of it
	indented := new(bytes.Buffer)
	err = json.Indent(indented, compact.Bytes(), "", "  ")
	if err != nil {
		return "", err
	}
	return indented.String(), nil
}
