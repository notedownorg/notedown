package tasks

import (
	"fmt"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
)

// file
func (c *Client) Create() {}

// file, line number
func (c *Client) Update() {}

func (c *Client) Delete(t ast.Task) error {
	err := c.writer.RemoveLine(writer.Document{Path: t.Path(), Hash: t.Version()}, t.Line())
	if err != nil {
		return fmt.Errorf("failed to remove task: %v: %w", t, err)
	}
	return nil
}

// new file, old file + line number
func (c *Client) Move() {}
