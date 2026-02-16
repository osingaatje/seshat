package context

//-----------------------------------//
// Context struct and helper methods //
//-----------------------------------//

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/MatusOllah/slogcolor"
	"github.com/fatih/color"
)

type Ctx struct {
	Queries Queries
	Trace   []Trace // records the order of queries

	Logger ContextLogger
}

type ContextLogger struct {
	log []string
	out io.Writer
}

func (l *ContextLogger) Write(b []byte) (n int, err error) {
	l.log = append(l.log, string(b))
	l.out.Write(b) // also write to 'out'

	return len(b), nil
}

func (l *ContextLogger) GetLogStrings() []string {
	return l.log
}

func (l *ContextLogger) GetLogString() string {
	return strings.Join(l.log, "")
}

type Trace struct {
	Funcname string
	Arg      any
}

func (c *Ctx) Error() string {
	return c.Logger.GetLogString()
}

func New() *Ctx {
	c := Ctx{}
	c.Queries = Queries{
		ctx: &c,
	}
	c.Trace = []Trace{}
	c.Logger = ContextLogger{log: []string{}, out: os.Stderr}

	// configure logging
	slog.SetDefault(slog.New(slogcolor.NewHandler(&c.Logger, logOpts())))

	return &c
}

func logOpts() *slogcolor.Options {
	opts := slogcolor.DefaultOptions
	opts.SrcFileMode = slogcolor.Nop
	opts.Level = slog.LevelDebug

	opts.LevelTags = map[slog.Level]string{
		slog.LevelDebug: color.New(color.BgBlack, color.Faint).Sprint("DEBUG"),
		slog.LevelInfo:  color.New(color.BgGreen, color.FgBlack).Sprint("INFO "),
		slog.LevelWarn:  color.New(color.BgYellow, color.FgBlack).Sprint("WARN "),
	}

	return opts
}

func (c *Ctx) LogDebug(msg string, args ...any) {
	slog.Debug(fmt.Sprintf(msg, args...))
}
func (c *Ctx) LogInfo(msg string, args ...any) {
	slog.Info(fmt.Sprintf(msg, args...))
}
func (c *Ctx) LogWarn(msg string, args ...any) {
	slog.Warn(fmt.Sprintf(msg, args...))
}
func (c *Ctx) LogErr(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
}
