package helper

import (
	"os"
)

func Export(filename string, object any) error {
	byt, err := MarshalJSONWithIndent(object)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, byt, os.ModePerm)
}
func ExportString(filename string, str string) error {
	return os.WriteFile(filename, []byte(str), os.ModePerm)
}
