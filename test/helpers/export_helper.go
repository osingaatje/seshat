package helpers

import (
	"os"

	"github.com/osingaatje/seshat/helper"
)

func Export(filename string, object any) error {
	byt, err := helper.MarshalJSONWithIndent(object)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, byt, os.ModePerm)
}
