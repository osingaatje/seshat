package shared

type GraphMetadata struct {
	Filename string `json:"filename"`
}

func NewGraphMetadata(filename string) GraphMetadata {
	return GraphMetadata{
		Filename: filename,
	}
}

func (m GraphMetadata) Copy() GraphMetadata {
	return GraphMetadata{
		Filename: m.Filename,
	}
}
