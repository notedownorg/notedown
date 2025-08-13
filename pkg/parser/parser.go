package parser

import (
	"bytes"

	"github.com/notedownorg/notedown/pkg/parser/extensions"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Parser defines the interface for parsing markdown documents
type Parser interface {
	Parse(source []byte) (*Document, error)
	ParseString(source string) (*Document, error)
}

// NotedownParser implements the Parser interface using goldmark
type NotedownParser struct {
	goldmark goldmark.Markdown
}

// NewParser creates a new Notedown parser
func NewParser() Parser {
	return &NotedownParser{
		goldmark: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				extensions.NewWikilinkExtension(),
			),
			goldmark.WithParserOptions(
				parser.WithAttribute(),
			),
		),
	}
}

// Parse parses markdown source bytes into a document tree
func (p *NotedownParser) Parse(source []byte) (*Document, error) {
	reader := text.NewReader(source)
	doc := p.goldmark.Parser().Parse(reader)

	return p.convertAST(doc, source), nil
}

// ParseString parses markdown source string into a document tree
func (p *NotedownParser) ParseString(source string) (*Document, error) {
	return p.Parse([]byte(source))
}

// convertAST converts goldmark AST to our custom tree structure
func (p *NotedownParser) convertAST(node ast.Node, source []byte) *Document {
	doc := NewDocument(Range{
		Start: Position{Line: 1, Column: 1, Offset: 0},
		End:   Position{Line: bytes.Count(source, []byte("\n")) + 1, Column: 1, Offset: len(source)},
	})

	p.convertNode(node, doc, source)
	return doc
}

// convertNode recursively converts goldmark AST nodes to our tree nodes
func (p *NotedownParser) convertNode(astNode ast.Node, parentNode Node, source []byte) {
	for child := astNode.FirstChild(); child != nil; child = child.NextSibling() {
		treeNode := p.astToTreeNode(child, source)
		if treeNode != nil {
			parentNode.AddChild(treeNode)

			// Only recurse for container nodes, not leaf nodes like headings
			switch child.(type) {
			case *ast.Heading:
				// Don't process heading children - text is already extracted
			case *ast.Text:
				// Text nodes are leaf nodes
			case *ast.CodeSpan:
				// Code spans are leaf nodes
			default:
				// For other nodes, process children
				p.convertNode(child, treeNode, source)
			}
		}
	}
}

// astToTreeNode converts a single goldmark AST node to our tree node
func (p *NotedownParser) astToTreeNode(astNode ast.Node, source []byte) Node {
	// Try to get position information
	var rng Range
	if hasText, ok := astNode.(interface{ Text([]byte) []byte }); ok {
		// For nodes with text, calculate range from text length
		txt := hasText.Text(source)
		rng = Range{
			Start: Position{Line: 1, Column: 1, Offset: 0}, // Default position
			End:   Position{Line: 1, Column: len(txt) + 1, Offset: len(txt)},
		}
	} else {
		// Default range for nodes without direct position info
		rng = Range{
			Start: Position{Line: 1, Column: 1, Offset: 0},
			End:   Position{Line: 1, Column: 1, Offset: 0},
		}
	}

	// Handle wikilink nodes
	if wikilink, ok := astNode.(*extensions.WikilinkAST); ok {
		return NewWikilink(wikilink.Target, wikilink.DisplayText, rng)
	}

	// Debug: Check for heading first
	if heading, ok := astNode.(*ast.Heading); ok {
		text := string(heading.Text(source))
		result := NewHeading(heading.Level, text, rng)
		return result
	}

	switch n := astNode.(type) {
	case *ast.Paragraph:
		return NewParagraph(rng)

	case *ast.Text:
		content := string(n.Segment.Value(source))
		return NewText(content, rng)

	case *ast.CodeBlock:
		var content bytes.Buffer
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			content.Write(line.Value(source))
		}

		return NewCodeBlock("", content.String(), false, rng)

	case *ast.FencedCodeBlock:
		var language string
		if n.Info != nil {
			info := n.Info.Text(source)
			if len(info) > 0 {
				language = string(info)
			}
		}

		var content bytes.Buffer
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			content.Write(line.Value(source))
		}

		return NewCodeBlock(language, content.String(), true, rng)

	case *ast.Link:
		url := string(n.Destination)
		var title string
		if n.Title != nil {
			title = string(n.Title)
		}
		return NewLink(url, title, rng)

	case *ast.List:
		return NewList(n.IsOrdered(), n.IsTight, rng)

	case *ast.ListItem:
		return NewListItem(false, false, rng) // TODO: Handle task lists

	case *ast.Emphasis:
		return NewEmphasis(rng)

	case *ast.CodeSpan:
		content := string(n.Text(source))
		return NewCode(content, rng)

	default:
		// Debug: Log unhandled node types
		// For unknown node types, create a generic container
		node := NewBaseNode(NodeContainer, rng)
		// If this was supposed to be a heading, fix it
		if heading, ok := astNode.(*ast.Heading); ok {
			text := string(heading.Text(source))
			result := NewHeading(heading.Level, text, rng)
			return result
		}
		return node
	}
}

// offsetToPosition converts byte offset to line/column position
func (p *NotedownParser) offsetToPosition(offset int, source []byte) Position {
	if offset > len(source) {
		offset = len(source)
	}

	line := 1
	column := 1

	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}
