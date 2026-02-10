package data

type DiagramFormat string

const (
	DiagramFormatUTML DiagramFormat = "utml" // extension name
)

type ParseCmd struct {
	DiagramFormat DiagramFormat
	Filepath      string
}

type ParseUTMLCmd struct {
	Filepath string
}
