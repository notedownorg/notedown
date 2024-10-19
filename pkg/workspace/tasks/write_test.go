package tasks_test

import (
	"fmt"
	"testing"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents/reader"
	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {

	client, _ := buildClient([]reader.Event{},
		// Create
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: "path"}, doc)
			assert.Equal(t, writer.AtEnd, line)
			return nil
		},

		// Update
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "update", method)
			assert.Equal(t, writer.Document{Path: "path", Hash: "version"}, doc)
			assert.Equal(t, 7, line)
			return nil
		},

		// Delete
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "remove", method)
			assert.Equal(t, writer.Document{Path: "path", Hash: "version"}, doc)
			assert.Equal(t, 1, line)
			return nil
		},

		// Move
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: "newPath"}, doc)
			assert.Equal(t, writer.AtEnd, line)
			return nil
		},
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "remove", method)
			assert.Equal(t, writer.Document{Path: "path", Hash: "version"}, doc)
			assert.Equal(t, 7, line)
			return nil
		},
	)

	assert.NoError(t, client.Create("path", "Task", ast.Todo, ast.WithLine(writer.AtEnd)))
	assert.NoError(t, client.Update(ast.NewTask(ast.NewIdentifier("path", "version"), "Task", ast.Todo, ast.WithLine(7))))
	assert.NoError(t, client.Delete(ast.NewTask(ast.NewIdentifier("path", "version"), "Task", ast.Todo, ast.WithLine(1))))
	assert.NoError(t, client.Move(ast.NewTask(ast.NewIdentifier("path", "version"), "Task", ast.Todo, ast.WithLine(7)), "newPath"))

}
