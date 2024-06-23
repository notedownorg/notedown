package ast

type Metadata map[string]interface{}

type Markers struct {
	ContentStart int
}

type Document struct {
	Metadata Metadata
	Markers  Markers
	Tasks    []Task
}
