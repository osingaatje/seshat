package generic

type GraphType string

const (
	UTMLResult  GraphType = "utml"
	ParseResult GraphType = "parse"
	Internal    GraphType = "internal"
	DotFile     GraphType = "dot" // you can view this at https://dreampuf.github.io/GraphvizOnline
)
