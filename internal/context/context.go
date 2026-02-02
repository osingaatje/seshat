package context

type Ctx struct {
	Queries Queries
	Trace   []Trace // records the order of queries

	// TODO error reporting etc.
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

	return &c
}

type Trace struct {
	Funcname string
	Arg      any
}
