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
	expDoc := writer.Document{Path: "path", Hash: "version"}

	client, _ := buildClient([]reader.Event{},
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "remove", method)
			assert.Equal(t, expDoc, doc)
			assert.Equal(t, 1, line)
			return nil
		},
	)

	assert.NoError(t, client.Delete(ast.NewTask(ast.NewIdentifier("path", "version"), "Task", ast.Todo, ast.WithLine(1))))

}
