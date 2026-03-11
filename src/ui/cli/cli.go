package cli

import (
	"fmt"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	g "github.com/osingaatje/seshat/types/generic"
	. "github.com/osingaatje/seshat/types/grade"

	"github.com/alecthomas/kong"
	dot "github.com/osingaatje/seshat/types/graph/dot"
	intern "github.com/osingaatje/seshat/types/graph/internal-rep"
	parse "github.com/osingaatje/seshat/types/graph/parse-result"
	utml "github.com/osingaatje/seshat/types/graph/utml"
)

// available options in the CLI
var cli struct {
	Test  TestCmd  `cmd:"" name:"test" help:"Print a test message using the Query system."`
	Emit  EmitCmd  `cmd:"" name:"emit" help:"Get some intermediate stage of the parsing process"`
	Grade GradeCmd `cmd:"" name:"grade" help:"Grade a solution"`
}

// main method
func RunCli(c *context.Ctx) {
	kongCli := kong.Parse(&cli, kong.Name("Seshat"))

	err := kongCli.Run(c)
	kongCli.FatalIfErrorf(err)
}

// ---------------- COMMANDS ---------------- //

// cmd EXAMPLE - command that uses the test query
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

// end EXAMPLE

// cmd EMIT
type EmitCmd struct {
	Input  string      `arg:"" required:"" name:"in" help:"Input UTML file" default:"./in.utml"`
	Output string      `arg:"" required:"" name:"out" help:"Output internal repr. file" default:"./out.json"`
	Type   g.GraphType `name:"type" help:"" default:"dot"`
}

func (cmd *EmitCmd) Run(c *context.Ctx) error {
	utml, parseresfixed, internalRep, dot, err := getReps(c, cmd.Input, cmd.Type)
	if err != nil {
		return err // logging is done internally in getReps
	}

	switch cmd.Type {
	case g.UTMLResult:
		return helper.Export(cmd.Output, utml)
	case g.ParseResult:
		return helper.Export(cmd.Output, parseresfixed)
	case g.Internal:
		return helper.Export(cmd.Output, internalRep)
	case g.DotFile:
		return helper.ExportString(cmd.Output, dot.String())
	}
	return logErrAndExit(c, "Unknown graph type %s", cmd.Type)
}

// end EMIT

// cmd GRADE
type GradeCmd struct {
	GradingScheme       string `name:"scheme" required:"" help:"Grading Scheme to use (JSON file)"`
	ReferenceSubmission string `arg:"" required:"" name:"reference" help:"Reference submission (.utml file)"`
	SubmissionDir       string `arg:"" required:"" name:"submission_directory" help:"Submission directory (containing .utml files)"`
}

func (cmd *GradeCmd) Run(c *context.Ctx) error {
	c.LogWarn("TODO GRADING RUBRIC!")
	rubric := GradeRubric{}

	allSubmissions, err := helper.AllUTMLFiles(cmd.SubmissionDir)
	if err != nil {
		return logErrAndExit(c, "%s", err.Error())
	}

	_, _, refRep, _, errRef := getReps(c, cmd.ReferenceSubmission, g.Internal)
	if errRef != nil {
		return logErrAndExit(c, "Could not parse Reference Submission: %s", errRef.Error())
	}

	results := map[string]*GradeResult{}

	for _, f := range allSubmissions {
		_, _, submissionRep, _, err := getReps(c, f, g.Internal)
		if err != nil {
			return logErrAndExit(c, "Could not parse student submission '%s': %s", f, err.Error())
		}

		res := c.Queries.GradeDiagram.Get("Grade diagram", command.GradeCmd{
			Rubric:            &rubric,
			ReferenceSolution: refRep,
			Submission:        submissionRep,
		})
		results[f] = res
	}

	return nil
}

// end GRADE

// ---------------- HELPERS ---------------- //
func logErrAndExit(c *context.Ctx, errMsg string, args ...any) error {
	c.LogErr(errMsg, args...)
	return fmt.Errorf(errMsg, args...)
}

func getReps(c *context.Ctx, inputFile string, graphType g.GraphType) (*utml.ParseResultUTML, *parse.ParseResult, *intern.InternalGraph, *dot.DotGraph, error) {
	// UTML
	utml := c.Queries.ParseUTML.Get("Parse UTML", inputFile)
	if utml == nil {
		return nil, nil, nil, nil, logErrAndExit(c, "No UTML parse result for '%s'", inputFile)
	}

	if graphType == g.UTMLResult {
		return utml, nil, nil, nil, nil
	}

	// PARSE RESULT
	parseres := c.Queries.ParseUTMLToInternal.Get("UTML -> Internal representation", utml)
	if parseres == nil {
		return utml, nil, nil, nil, logErrAndExit(c, "Failed conversion to internal representation.")
	}

	// FIXED PARSE RESULT
	fixed := c.Queries.RepairDiagram.Get("Repair internal repr.", command.NewRepairCmdDefOpt(
		parseres,
	))
	if fixed == nil {
		return utml, parseres, nil, nil, logErrAndExit(c, "Failed fixing diagram '%s'!", inputFile)
	}
	if graphType == g.ParseResult {
		return utml, fixed, nil, nil, nil
	}

	if graphType == g.DotFile {
		dot := c.Queries.DisplayDiagramAsDot.Get("Internal -> .dot", fixed)
		if dot == nil {
			return utml, fixed, nil, nil, logErrAndExit(c, "Could not convert '%s' from internal to .dot file", inputFile)
		}

		return utml, fixed, nil, dot, nil
	}

	if graphType == g.Internal {
		intern := c.Queries.ConvertGraphToInternal.Get("Parseres -> Internal", fixed)
		if intern == nil {
			return utml, fixed, nil, nil, logErrAndExit(c, "Could not convert Parse Result to internal for file '%s'", inputFile)
		}
		return utml, fixed, intern, nil, nil
	}

	return nil, nil, nil, nil, logErrAndExit(c, "Unknown graphtype requested: %s", graphType)
}
