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
	prefixes []string // for logging function name while in that function etc.
	log      []string
	out      io.Writer
}

func NewContextLogger() ContextLogger {
	return ContextLogger{
		prefixes: []string{},
		log:      []string{},
		out:      os.Stderr,
	}
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
	c.Logger = NewContextLogger()

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

const MAX_PREFIX_CHARS = 50

func (c *Ctx) LogPrefixAdd(pref string, args ...any) {
	pr := fmt.Sprintf(pref, args...)
	if len(pr) > MAX_PREFIX_CHARS {
		pr = pr[:MAX_PREFIX_CHARS]
	}
	c.Logger.prefixes = append(c.Logger.prefixes, pr)
}
func (c *Ctx) LogPrefixRm() {
	prefLen := len(c.Logger.prefixes)
	if prefLen == 0 {
		return
	}
	c.Logger.prefixes = c.Logger.prefixes[:prefLen-1]
}

func (c *Ctx) LogDebug(msg string, args ...any) {
	slog.Debug(c.formatLogMsg(msg, args...))
}
func (c *Ctx) LogInfo(msg string, args ...any) {
	slog.Info(c.formatLogMsg(msg, args...))
}
func (c *Ctx) LogWarn(msg string, args ...any) {
	slog.Warn(c.formatLogMsg(msg, args...))
}
func (c *Ctx) LogErr(msg string, args ...any) {
	slog.Error(c.formatLogMsg(msg, args...))
}
func (c *Ctx) LogErrAndReturn(errMsg string, args ...any) error {
	c.LogErr(errMsg, args...)
	return fmt.Errorf(errMsg, args...)
}

func (c *Ctx) formatLogMsg(msg string, args ...any) string {
	rawLogMsg := fmt.Sprintf(msg, args...)
	if len(c.Logger.prefixes) == 0 {
		return rawLogMsg
	}
	return strings.Join(c.Logger.prefixes, " - ") + ": " + rawLogMsg
}
