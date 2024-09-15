package ast

type Metadata map[string]interface{}

type Markers struct {
	ContentStart int
}

type Document struct {
	Hash     string
	Metadata Metadata
	Markers  Markers
	Tasks    []Task
}
