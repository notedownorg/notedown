package writer

import "fmt"

type LineWriter interface {
	AddLine(doc Document, line int, obj fmt.Stringer) error
	RemoveLine(doc Document, line int) error
	UpdateLine(doc Document, line int, obj fmt.Stringer) error
}
