package test

//-----------------//
// An example test //
//-----------------//

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/osingaatje/seshat/context/data"
	"github.com/osingaatje/seshat/driver"
)

func TestTestCmd(t *testing.T) {
	c := driver.NewContext()

	res := c.Queries.Test.Get("Test", data.NameCmd{FName: "Douwe", LName: "O"})
	assert.Equal(t, "Hello Douwe O!", res)
}
