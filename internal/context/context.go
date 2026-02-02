package context

type Ctx struct {
	trace []Trace
}

type Trace struct {
	Funcname string
	Arg      any
}
