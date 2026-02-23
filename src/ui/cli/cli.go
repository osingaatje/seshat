package cli

import (
	"fmt"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types"

	"github.com/alecthomas/kong"
)

// Example command that uses the test query
type TestCmd struct {
	FName string `arg:"" name:"fname" help:"Your first name :)"`
	LName string `arg:"" name:"lname" help:"Your last name :)"`
}

// method that binds to the test command
func (t *TestCmd) Run(c *context.Ctx) error {
	c.LogDebug("test 123 debug")
	c.LogInfo("test 123 info")
	c.LogWarn("test 123 warn")
	c.LogErr("test 123 error")

	fmt.Println(c.Queries.Test.Get("Test", types.NameCmd{FName: t.FName, LName: t.LName}))
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
