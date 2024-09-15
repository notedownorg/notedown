package documents

import (
	"github.com/liamawhite/nl/pkg/ast"
)

type Document struct {
	lastUpdated int64
	Hash        string
	ast.Document
}
