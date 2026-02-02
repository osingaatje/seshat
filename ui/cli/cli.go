package cli

import (
	"fmt"
	"seshat/context"

	"github.com/alecthomas/kong"
)

var cli struct {
	Test TestCmd `cmd:"" name:"test" help:"Print a test message using the Query system."`
}

// internal cli commands
// -- TEST -- //
type TestCmd struct {
	Name string `arg:"" name:"name" help:"Your name :)"`
}

func (t *TestCmd) Run(c *context.Ctx) error {
	fmt.Println(c.Queries.Test.Get("Test", t.Name))
	return nil
}

// --------- //

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}
