package parser

// Visitor defines the visitor pattern interface for traversing the document tree
type Visitor interface {
	Visit(node Node) error
}

// Walker provides utilities for traversing document trees
type Walker struct {
	visitor Visitor
}

// NewWalker creates a new tree walker
func NewWalker(visitor Visitor) *Walker {
	return &Walker{visitor: visitor}
}

// Walk traverses the tree starting from the given node
func (w *Walker) Walk(node Node) error {
	if err := node.Accept(w.visitor); err != nil {
		return err
	}

	for _, child := range node.Children() {
		if err := w.Walk(child); err != nil {
			return err
		}
	}

	return nil
}

// WalkFunc is a convenience type for creating visitors from functions
type WalkFunc func(node Node) error

// Visit implements the Visitor interface
func (f WalkFunc) Visit(node Node) error {
	return f(node)
}
