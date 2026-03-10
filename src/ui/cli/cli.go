package cli

import (
	"fmt"

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

type EmitCmd struct {
	Input  string    `arg:"" required:"" name:"in" help:"Input UTML file" default:"./in.utml"`
	Output string    `arg:"" required:"" name:"out" help:"Output internal repr. file" default:"./out.json"`
	Type   GraphType `name:"type" help:"" default:"dot"`
}
type GraphType string

const (
	ParseResult GraphType = "parse"
	DotFile     GraphType = "dot" // you can view this at https://dreampuf.github.io/GraphvizOnline
	UTMLResult  GraphType = "utml"
	Internal    GraphType = "internal"
)

func (cmd *EmitCmd) Run(c *context.Ctx) error {
	// UTML
	utml := c.Queries.ParseUTML.Get("Parse UTML", command.ParseUTMLCmd{Filepath: cmd.Input})
	if utml == nil {
		return logErrAndExit(c, "No UTML parse result for '%s'", cmd.Input)
	}

	if cmd.Type == UTMLResult {
		return helper.Export(cmd.Output, utml)
	}

	// PARSE RESULT
	parseres := c.Queries.ParseUTMLToInternal.Get("UTML -> Internal representation", utml)
	if parseres == nil {
		return logErrAndExit(c, "Failed conversion to internal representation.")
	}

	// FIXED PARSE RESULT
	fixed := c.Queries.RepairDiagram.Get("Repair internal repr.", command.NewRepairCmdDefOpt(
		parseres,
	))
	if fixed == nil {
		return logErrAndExit(c, "Failed fixing diagram '%s'!", cmd.Input)
	}
	if cmd.Type == ParseResult {
		return helper.Export(cmd.Output, parseres)
	}

	if cmd.Type == DotFile {
		dot := c.Queries.DisplayDiagramAsDot.Get("Internal -> .dot", fixed)
		if dot == nil {
			return logErrAndExit(c, "Could not convert '%s' from internal to .dot file", cmd.Input)
		}

		return helper.ExportString(cmd.Output, dot.String())
	}

	if cmd.Type == Internal {
		intern := c.Queries.ConvertGraphToInternal.Get("Parseres -> Internal", fixed)
		if intern == nil {
			return logErrAndExit(c, "Could not convert Parse Result to internal for file '%s'", cmd.Input)
		}
		return helper.Export(cmd.Output, intern)
	}

	return nil
}

// end INTERNAL REPR

// available options in the CLI
var cli struct {
	Test TestCmd `cmd:"" name:"test" help:"Print a test message using the Query system."`
	Emit EmitCmd `cmd:"" name:"emit" help:"Get some intermediate stage of the parsing process"`
}

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}

func logErrAndExit(c *context.Ctx, errMsg string, args ...any) error {
	c.LogErr(errMsg, args...)
	return fmt.Errorf(errMsg, args...)
}
