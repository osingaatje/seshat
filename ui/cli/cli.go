package cli

import (
	"fmt"
	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/context/data"

	"github.com/alecthomas/kong"
)

// Example command that uses the test query
type TestCmd struct {
	FName string `arg:"" name:"fname" help:"Your first name :)"`
	LName string `arg:"" name:"lname" help:"Your last name :)"`
}

// method that binds to the test command
func (t *TestCmd) Run(c *context.Ctx) error {
	fmt.Println(c.Queries.Test.Get("Test", data.NameCmd{FName: t.FName, LName: t.LName}))
	return nil
}

// available options in the CLI
var cli struct {
	Test TestCmd `cmd:"" name:"test" help:"Print a test message using the Query system."`
}

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}
