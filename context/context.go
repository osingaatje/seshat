package context

import (
	"log/slog"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/fatih/color"
)

type Ctx struct {
	Queries Queries
	Trace   []Trace // records the order of queries

	// TODO perhaps centrally collect logs? meh.
}

type Trace struct {
	Funcname string
	Arg      any
}

func (c *Ctx) Error() string {
	return "TODO error ctx"
}

func New() *Ctx {
	c := Ctx{}
	c.Queries = Queries{
		ctx: &c,
	}
	c.Trace = []Trace{}

	// configure logging
	opts := slogcolor.DefaultOptions
	opts.SrcFileMode = slogcolor.Nop
	opts.Level = slog.LevelDebug

	opts.LevelTags = map[slog.Level]string{
		slog.LevelDebug: color.New(color.BgBlack, color.Faint).Sprint("DEBUG"),
		slog.LevelInfo:  color.New(color.BgGreen, color.FgBlack).Sprint("INFO"),
		slog.LevelWarn:  color.New(color.BgYellow, color.FgBlack).Sprint("WARN"),
	}

	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, opts)))

	return &c
}

func (c *Ctx) LogDebug(msg string, args ...any) {
	slog.Debug(msg, args...)
}
func (c *Ctx) LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}
func (c *Ctx) LogWarning(msg string, args ...any) {
	slog.Warn(msg, args...)
}
func (c *Ctx) LogError(msg string, args ...any) {
	slog.Error(msg, args...)
}
