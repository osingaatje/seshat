package cli

import (
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

// INTERNAL REPR
type InternalRepCmd struct {
	Input  string `arg:"" name:"in" help:"Input UTML file"`
	Output string `arg:"" name:"out" help:"Output internal repr. file" default:"./out.json"`
}

func (cmd *InternalRepCmd) Run(c *context.Ctx) error {
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: cmd.Input})
	if utml == nil {
		return fmt.Errorf("No UTML parse result for '%s'", cmd.Input)
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

// end INTERNAL REPR

// PRINT GRAPH .dot FILE
type PrintGraphCmd struct {
	Input  string `arg:"" name:"in" help:"Input file (.utml)"`
	Output string `arg:"" name:"out" help:"Output file (.dot)"`
}

func (cmd *PrintGraphCmd) Run(c *context.Ctx) error {
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: cmd.Input})
	if utml == nil {
		return fmt.Errorf("No UTML parse res for file '%s'", cmd.Input)
	}
	intern := c.Queries.ParseUTMLToInternal.Get("UTML -> Internal", utml)
	if intern == nil {
		return fmt.Errorf("Could not convert '%s' to internal result", cmd.Input)
	}

	dot := c.Queries.DisplayInternalReprToDot.Get("Internal -> .dot", intern)
	if dot == nil {
		return fmt.Errorf("Could not convert '%s' from internal to .dot file", cmd.Input)
	}

	err := os.WriteFile(cmd.Output, []byte(dot.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not write to '%s', err=%s", cmd.Output, err.Error())
	}
	return nil
}

// available options in the CLI
var cli struct {
	Test           TestCmd        `cmd:"" name:"test" help:"Print a test message using the Query system."`
	GetInternalRep InternalRepCmd `cmd:"" name:"get-repr" help:"Print the internal representation for some UTML file"`
	PrintDot       PrintGraphCmd  `cmd:"" name:"get-dot" help:"Export .dot file from .utml"`
}

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}
