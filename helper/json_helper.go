package helper

import (
	"bytes"
	"encoding/json"
	"io"
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

// Explicitly does not escape HTML codes. That messes up associations containing '>' for ex.
func MarshalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	writer := io.Writer(&buf)

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false) // this is why we need a custom encoder. Stupid automatic HTML escaping.

	err := enc.Encode(v)
	return buf.Bytes(), err
}
