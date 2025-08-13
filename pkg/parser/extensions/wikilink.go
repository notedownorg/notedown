package extensions

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// WikilinkExtension adds support for wikilinks ([[page]] and [[page|display]])
type WikilinkExtension struct{}

// Extend implements goldmark.Extender
func (e *WikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&wikilinkParser{}, 200),
	))
}

// NewWikilinkExtension creates a new wikilink extension
func NewWikilinkExtension() goldmark.Extender {
	return &WikilinkExtension{}
}

// wikilinkParser parses wikilink syntax
type wikilinkParser struct{}

var wikilinkRegex = regexp.MustCompile(`^\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

// Trigger returns the trigger characters for wikilinks
func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

// Parse parses a wikilink
func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 4 || line[0] != '[' || line[1] != '[' {
		return nil
	}

	match := wikilinkRegex.FindSubmatch(line)
	if match == nil {
		return nil
	}

	// Extract target and display text
	target := strings.TrimSpace(string(match[1]))
	var displayText string
	if len(match) > 2 && match[2] != nil {
		displayText = strings.TrimSpace(string(match[2]))
	} else {
		displayText = target
	}

	// Create wikilink AST node
	node := &WikilinkAST{
		Target:      target,
		DisplayText: displayText,
	}

	// Note: For now, we'll rely on goldmark's internal segment handling

	// Advance the reader
	block.Advance(len(match[0]))

	return node
}

// WikilinkAST represents a wikilink in the goldmark AST
type WikilinkAST struct {
	ast.BaseInline
	Target      string
	DisplayText string
}

// Dump implements ast.Node
func (n *WikilinkAST) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Target":      n.Target,
		"DisplayText": n.DisplayText,
	}, nil)
}

// Kind returns the node kind
func (n *WikilinkAST) Kind() ast.NodeKind {
	return WikilinkKind
}

// WikilinkKind is the kind for wikilink nodes
var WikilinkKind = ast.NewNodeKind("Wikilink")
