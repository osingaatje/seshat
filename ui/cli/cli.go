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
	FName string `arg:"" name:"fname" help:"Your first name :)"`
	LName string `arg:"" name:"lname" help:"Your last name :)"`
}

func (t *TestCmd) Run(c *context.Ctx) error {
	c.LogDebug("TEST!")
	c.LogInfo("info")
	c.LogWarning("warning!")
	c.LogError("errror1!!1!!!!")
	fmt.Println(c.Queries.Test.Get("Test", context.MultiKey2[string, string]{K1: t.FName, K2: t.LName}))
	return nil
}

// --------- //

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}
