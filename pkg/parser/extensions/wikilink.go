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
		util.Prioritized(&wikilinkParser{}, 100), // Higher priority than link (200)
	))
}

// NewWikilinkExtension creates a new wikilink extension
func NewWikilinkExtension() goldmark.Extender {
	return &WikilinkExtension{}
}

// wikilinkParser parses wikilink syntax
type wikilinkParser struct{}

var wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

// Trigger returns the trigger characters for wikilinks
func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

// Parse parses a wikilink
func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	if len(line) < 4 {
		return nil
	}

	if line[0] != '[' || line[1] != '[' {
		return nil
	}

	// Find the closing ]]
	closePos := -1
	for i := 2; i < len(line)-1; i++ {
		if line[i] == ']' && line[i+1] == ']' {
			closePos = i
			break
		}
	}

	if closePos == -1 {
		return nil
	}

	// Extract the content between [[ and ]]
	content := line[2:closePos]
	wikilinkLength := closePos + 2

	// Parse target and display text
	target := string(content)
	displayText := target

	// Check for pipe separator
	if pipePos := strings.Index(target, "|"); pipePos != -1 {
		displayText = strings.TrimSpace(target[pipePos+1:])
		target = strings.TrimSpace(target[:pipePos])
	}
	target = strings.TrimSpace(target)
	displayText = strings.TrimSpace(displayText)

	if target == "" {
		return nil
	}

	// Reject targets containing .. sequences to prevent directory traversal
	normalizedTarget := strings.ReplaceAll(target, "\\", "/")
	if strings.Contains(normalizedTarget, "..") {
		return nil
	}

	// Create wikilink AST node
	node := &WikilinkAST{
		Target:      target,
		DisplayText: displayText,
	}

	// Advance the reader by the length of the wikilink
	block.Advance(wikilinkLength)

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
