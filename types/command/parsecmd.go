package command

type DiagramFormat int

const (
	DiagramFormatUTML DiagramFormat = iota
)

type ParseCmd struct {
	DiagramFormat DiagramFormat
	Filepath      string
}

type ParseUTMLCmd string
