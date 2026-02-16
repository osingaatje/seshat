package data

// go:exhaustive
type DiagramFormat int

const (
	DiagramFormatUTML DiagramFormat = iota
)

type ParseCmd struct {
	DiagramFormat DiagramFormat
	Filepath      string
}

type ParseUTMLCmd struct {
	Filepath string
}
