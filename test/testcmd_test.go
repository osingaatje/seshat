package test

//-----------------//
// An example test //
//-----------------//

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/types"
)

func TestTestCmd(t *testing.T) {
	c := driver.NewContext()

	res := c.Queries.Test.Get("Test", data.NameCmd{FName: "Douwe", LName: "O"})
	assert.Equal(t, "Hello Douwe O!", res)
}
