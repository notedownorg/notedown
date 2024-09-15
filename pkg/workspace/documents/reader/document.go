package reader

import (
	"github.com/liamawhite/nl/pkg/ast"
)

type Document struct {
	lastUpdated int64
	ast.Document
}
