package ast

type Metadata map[string]interface{}

type Document struct {
	Metadata Metadata
	Tasks    []Task
}
