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

var ctx = driver.NewContext()

func TestTestCmd(t *testing.T) {
	res := ctx.Queries.Test.Get("Test", data.NameCmd{FName: "Douwe", LName: "O"})
	assert.Equal(t, "Hello Douwe O!", res)
}
