package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osingaatje/seshat/helper"
	"github.com/osingaatje/seshat/src/context"
	"github.com/osingaatje/seshat/types/command"
	g "github.com/osingaatje/seshat/types/generic"
	"github.com/osingaatje/seshat/types/grade"

	"github.com/alecthomas/kong"
	displaygraph "github.com/osingaatje/seshat/types/graph/dot"
	dot "github.com/osingaatje/seshat/types/graph/dot"
	parse "github.com/osingaatje/seshat/types/graph/intern"
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
	Input  string      `required:"" name:"in" help:"Input UTML file(s)"`
	Output string      `required:"" name:"out" help:"Output directory"`
	Type   g.GraphType `name:"type" help:"" default:"dot"`
}

func (cmd *EmitCmd) Run(c *context.Ctx) error {
	names, utml, internal, dot, err := getRepsOfFiles(c, cmd.Input, cmd.Type)
	if err != nil {
		return err // logging is done internally in getReps
	}

	switch cmd.Type {
	case g.UTMLResult:
		return exportMany(cmd.Output, names, utml)

	case g.InternalRep:
		return exportMany(cmd.Output, names, internal)

	case g.DotFile:
		if dot == nil {
			return c.LogErrAndReturn("Dot graphs were nil! BUG IN CODE")
		}
		dotstrings := helper.Map(dot, func(d *displaygraph.DotGraph) string { return d.String() })
		return exportManyString(cmd.Output, names, dotstrings)
	}

	return c.LogErrAndReturn("Unknown graph type %s", cmd.Type)
}

func outFile(dir, infile string) string {
	return filepath.Join(dir, filepath.Base(infile)+"-fixed.json")
}

func exportMany[E any](basedir string, names []string, objs []E) error {
	return exportManyFunc(basedir, names, objs, func(fn string, obj E) error { return helper.Export(fn, obj) })
}

func exportManyString(basedir string, names []string, objs []string) error {
	return exportManyFunc(basedir, names, objs, func(fn string, obj string) error { return helper.ExportString(fn, obj) })
}

func exportManyFunc[E any](basedir string, names []string, objs []E, exportFunc func(filename string, obj E) error) error {
	if len(names) != len(objs) {
		return fmt.Errorf("BUG IN TOOL: NAMES LENGTH != EXPORT OBJECT LENGTH")
	}

	for i, res := range objs {
		var e error
		filename := outFile(basedir, names[i])

		e = exportFunc(filename, res)
		if e != nil {
			return e
		}
	}
	return nil
}

// end EMIT

// cmd GRADE
type GradeCmd struct {
	GradingScheme       string `name:"scheme" required:"" help:"Grading Scheme to use (JSON file)"`
	ReferenceSubmission string `name:"solution" help:"Reference submission (.utml file)"`
	SubmissionDir       string `name:"submission-dir" help:"Submission directory (containing .utml files)"`
	ResultsDir          string `name:"results-dir" help:"Which directory to write resulting files to"`
}

func (cmd *GradeCmd) Run(c *context.Ctx) error {
	rubric := grade.GradeRubric{}
	dat, err := os.ReadFile(cmd.GradingScheme)
	if err != nil {
		return c.LogErrAndReturn("Could not load grading scheme from path '%s': %s", cmd.GradingScheme, err.Error())
	}
	err = helper.UnmarshalJSON(dat, &rubric)
	if err != nil {
		return c.LogErrAndReturn("Incorrect grading rubric at '%s': %s", cmd.GradingScheme, err.Error())
	}

	allSubmissions, err := helper.AllUTMLFiles(cmd.SubmissionDir)
	if err != nil {
		return c.LogErrAndReturn("Error loading UTML files from directory '%s': %s", cmd.SubmissionDir, err.Error())
	}

	_, refRep, _, err := getReps(c, cmd.ReferenceSubmission, g.InternalRep)
	if err != nil {
		return c.LogErrAndReturn("Could not parse Reference Submission: %s", err.Error())
	}

	err = os.MkdirAll(cmd.ResultsDir, os.ModePerm)
	if err != nil {
		return c.LogErrAndReturn("Could not verify/make results directory")
	}

	var submissionErr error = nil
	for _, f := range allSubmissions {
		_, submissionRep, _, err := getReps(c, f, g.InternalRep)
		if err != nil {
			return c.LogErrAndReturn("Could not parse student submission '%s': %s", f, err.Error())
		}

		res := c.Queries.GradeDiagram.Get("Grade diagram", command.GradeCmd{
			Rubric:            &rubric,
			ReferenceSolution: refRep,
			Submission:        submissionRep,
		})

		if res == nil {
			submissionErr = c.LogErrAndReturn("Could not grade diagram '%s'!", f)
			continue
		}
		jsonBytes, err := helper.MarshalJSON(res)
		if err != nil {
			submissionErr = c.LogErrAndReturn("Could not marshal '%s' to JSON: %s", f, err.Error())
			continue
		}
		gradingFilename := filepath.Join(cmd.ResultsDir, f+"-grade.json")

		err = os.WriteFile(gradingFilename, jsonBytes, os.ModePerm)
		if err != nil {
			submissionErr = c.LogErrAndReturn("Could not write grading config '%s': %s", gradingFilename, err.Error())
			continue
		}
	}

	return submissionErr
}

// end GRADE

// ---------------- HELPERS ---------------- //

func getRepsOfFiles(c *context.Ctx, inputGlob string, graphType g.GraphType) ([]string, []*utml.ParseResultUTML, []*parse.InternalGraph, []*dot.DotGraph, error) {
	matches, err := filepath.Glob(inputGlob)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if len(matches) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("No matches found for GLOB '%s'", inputGlob)
	}
	names := []string{}
	resU, resP, resD := []*utml.ParseResultUTML{}, []*parse.InternalGraph{}, []*dot.DotGraph{}
	for _, match := range matches {
		names = append(names, match)
		u, p, d, e := getReps(c, match, graphType)
		if e != nil {
			return names, resU, resP, resD, e
		}

		resU = append(resU, u)
		resP = append(resP, p)
		resD = append(resD, d)
	}
	return names, resU, resP, resD, nil
}

func getReps(c *context.Ctx, inputFile string, graphType g.GraphType) (*utml.ParseResultUTML, *parse.InternalGraph, *dot.DotGraph, error) {
	// UTML
	utml := c.Queries.ParseUTML.Get("Parse UTML", inputFile)
	if utml == nil {
		return nil, nil, nil, c.LogErrAndReturn("No UTML parse result for '%s'", inputFile)
	}

	if graphType == g.UTMLResult {
		return utml, nil, nil, nil
	}

	// PARSE RESULT
	parseres := c.Queries.ParseUTMLToParseRes.Get("UTML -> Internal representation", utml)
	if parseres == nil {
		return utml, nil, nil, c.LogErrAndReturn("Failed conversion to internal representation.")
	}

	// FIXED PARSE RESULT
	repairRes := c.Queries.RepairDiagram.Get("Repair internal repr.", command.NewRepairCmdDefOpt(
		parseres,
	))
	if len(repairRes.Errors) > 0 || repairRes.Diagram == nil {
		return utml, parseres, nil, c.LogErrAndReturn("Failed fixing diagram '%s'!", inputFile)
	}
	fixedDiag := repairRes.Diagram
	if graphType == g.InternalRep {
		return utml, fixedDiag, nil, nil
	}

	if graphType == g.DotFile {
		dot := c.Queries.DisplayDiagramAsDot.Get("Internal -> .dot", fixedDiag)
		if dot == nil {
			return utml, fixedDiag, nil, c.LogErrAndReturn("Could not convert '%s' from internal to .dot file", inputFile)
		}

		return utml, fixedDiag, dot, nil
	}

	return nil, nil, nil, c.LogErrAndReturn("Unknown graphtype requested: %s", graphType)
}
