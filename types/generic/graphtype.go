package generic

type GraphType string

const (
	UTMLResult  GraphType = "utml"
	InternalRep GraphType = "parse"
	DotFile     GraphType = "dot" // you can view this at https://dreampuf.github.io/GraphvizOnline
)
