package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"

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

	fmt.Println(c.Queries.Test.Get("Test", command.NameCmd{FName: t.FName, LName: t.LName}))
	return nil
}

type InternalRepCmd struct {
	Input  string `arg:"" name:"in" help:"Input UTML file"`
	Output string `arg:"" name:"out" help:"Output internal repr. file" default:"./out.json"`
}

func (cmd *InternalRepCmd) Run(c *context.Ctx) error {
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: cmd.Input})
	if utml == nil {
		errMsg := fmt.Sprintf("No UTML parse result for '%s'", cmd.Input)
		c.LogErr("%s", errMsg)
		return errors.New(errMsg)
	}

	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> Internal representation", utml)
	if intern == nil {
		c.LogWarn("Internal representation is nil!")
		return nil
	}

	intern_bytes, err := helper.MarshalJSONWithIndent(intern)
	if err != nil {
		c.LogErr("Could not marshal internal repr. to JSON")
		return err
	}

	err = os.WriteFile(cmd.Output, intern_bytes, os.ModePerm)
	if err != nil {
		c.LogErr("Could not write to file :( err=%s", err.Error())
		return err
	}
	return nil
}

// available options in the CLI
var cli struct {
	Test           TestCmd        `cmd:"" name:"test" help:"Print a test message using the Query system."`
	GetInternalRep InternalRepCmd `cmd:"" name:"get-repr" help:"Print the internal representation for some UTML file"`
}

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}
