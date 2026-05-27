package helper

import (
	"bytes"
	"encoding/json"
	"io"
)

const INDENT_PREFIX = ""
const INDENT_STR = "  "

// Remarshals JSON, ensuring sorted keys
func RemarshalJSON(bytes []byte) ([]byte, error) {
	var jsObj any
	err := json.Unmarshal(bytes, &jsObj)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsObj)
}

func IndentJSON(fileContents []byte) ([]byte, error) {
	compact := new(bytes.Buffer)
	// first compact it
	err := json.Compact(compact, fileContents)
	if err != nil {
		return nil, err
	}

	// now indent the shit out of it
	indented := new(bytes.Buffer)
	err = json.Indent(indented, compact.Bytes(), INDENT_PREFIX, INDENT_STR)
	if err != nil {
		return nil, err
	}
	return indented.Bytes(), nil
}

// Explicitly does not escape HTML codes. That messes up associations containing '>' for ex.
func MarshalJSON(v any) ([]byte, error) {
	return marshJson(v, false)
}

func MarshalJSONWithIndent(v any) ([]byte, error) {
	return marshJson(v, true)
}

func marshJson(v any, setIndent bool) ([]byte, error) {
	var buf bytes.Buffer
	writer := io.Writer(&buf)

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false) // this is why we need a custom encoder. Stupid automatic HTML escaping.

	if setIndent {
		enc.SetIndent(INDENT_PREFIX, INDENT_STR)
	}

	err := enc.Encode(v)
	return buf.Bytes(), err
}

// stores result in "res" which needs to be a POINTER!, gives error "error" if error occurred.
func UnmarshalJSON(data []byte, res any) error {
	reader := bytes.NewReader(data)

	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields() // FORCE ERROR ON ANY UNKNOWN FIELDS!

	return dec.Decode(res)
}
