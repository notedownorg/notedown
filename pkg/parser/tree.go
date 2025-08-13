package parser

import "fmt"

// Position represents a position in the source document
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // 0-based byte offset
}

// Range represents a range in the source document
type Range struct {
	Start Position
	End   Position
}

// NodeType represents the type of a node in the document tree
type NodeType int

const (
	// Block nodes
	NodeDocument NodeType = iota
	NodeHeading
	NodeParagraph
	NodeCodeBlock
	NodeBlockQuote
	NodeList
	NodeListItem
	NodeThematicBreak

	// Inline nodes
	NodeText
	NodeEmphasis
	NodeStrong
	NodeCode
	NodeLink
	NodeWikilink
	NodeAutoLink
	NodeRawHTML

	// Container nodes
	NodeContainer
)

// String returns the string representation of a NodeType
func (nt NodeType) String() string {
	switch nt {
	case NodeDocument:
		return "Document"
	case NodeHeading:
		return "Heading"
	case NodeParagraph:
		return "Paragraph"
	case NodeCodeBlock:
		return "CodeBlock"
	case NodeBlockQuote:
		return "BlockQuote"
	case NodeList:
		return "List"
	case NodeListItem:
		return "ListItem"
	case NodeThematicBreak:
		return "ThematicBreak"
	case NodeText:
		return "Text"
	case NodeEmphasis:
		return "Emphasis"
	case NodeStrong:
		return "Strong"
	case NodeCode:
		return "Code"
	case NodeLink:
		return "Link"
	case NodeWikilink:
		return "Wikilink"
	case NodeAutoLink:
		return "AutoLink"
	case NodeRawHTML:
		return "RawHTML"
	case NodeContainer:
		return "Container"
	default:
		return fmt.Sprintf("Unknown(%d)", int(nt))
	}
}

// Node represents a node in the document tree
type Node interface {
	Type() NodeType
	Range() Range
	Parent() Node
	SetParent(Node)
	Children() []Node
	AddChild(Node)
	RemoveChild(Node)
	Accept(Visitor) error
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	kind     NodeType
	range_   Range
	parent   Node
	children []Node
}

// NewBaseNode creates a new base node
func NewBaseNode(kind NodeType, rng Range) *BaseNode {
	return &BaseNode{
		kind:     kind,
		range_:   rng,
		children: make([]Node, 0),
	}
}

// Type returns the node type
func (b *BaseNode) Type() NodeType {
	return b.kind
}

// Range returns the source range
func (b *BaseNode) Range() Range {
	return b.range_
}

// Parent returns the parent node
func (b *BaseNode) Parent() Node {
	return b.parent
}

// SetParent sets the parent node
func (b *BaseNode) SetParent(parent Node) {
	b.parent = parent
}

// Children returns the child nodes
func (b *BaseNode) Children() []Node {
	return b.children
}

// AddChild adds a child node
func (b *BaseNode) AddChild(child Node) {
	child.SetParent(b)
	b.children = append(b.children, child)
}

// RemoveChild removes a child node
func (b *BaseNode) RemoveChild(child Node) {
	for i, c := range b.children {
		if c == child {
			child.SetParent(nil)
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
}

// Accept implements the visitor pattern
func (b *BaseNode) Accept(visitor Visitor) error {
	return visitor.Visit(b)
}

// Document represents the root document node
type Document struct {
	*BaseNode
	Title string
}

// NewDocument creates a new document node
func NewDocument(rng Range) *Document {
	return &Document{
		BaseNode: NewBaseNode(NodeDocument, rng),
	}
}

// Heading represents a heading node
type Heading struct {
	*BaseNode
	Level int
	Text  string
}

// NewHeading creates a new heading node
func NewHeading(level int, text string, rng Range) *Heading {
	return &Heading{
		BaseNode: NewBaseNode(NodeHeading, rng),
		Level:    level,
		Text:     text,
	}
}

// Paragraph represents a paragraph node
type Paragraph struct {
	*BaseNode
}

// NewParagraph creates a new paragraph node
func NewParagraph(rng Range) *Paragraph {
	return &Paragraph{
		BaseNode: NewBaseNode(NodeParagraph, rng),
	}
}

// Text represents a text node
type Text struct {
	*BaseNode
	Content string
}

// NewText creates a new text node
func NewText(content string, rng Range) *Text {
	return &Text{
		BaseNode: NewBaseNode(NodeText, rng),
		Content:  content,
	}
}

// CodeBlock represents a code block node
type CodeBlock struct {
	*BaseNode
	Language string
	Content  string
	Fenced   bool
}

// NewCodeBlock creates a new code block node
func NewCodeBlock(language, content string, fenced bool, rng Range) *CodeBlock {
	return &CodeBlock{
		BaseNode: NewBaseNode(NodeCodeBlock, rng),
		Language: language,
		Content:  content,
		Fenced:   fenced,
	}
}

// Link represents a link node
type Link struct {
	*BaseNode
	URL   string
	Title string
}

// NewLink creates a new link node
func NewLink(url, title string, rng Range) *Link {
	return &Link{
		BaseNode: NewBaseNode(NodeLink, rng),
		URL:      url,
		Title:    title,
	}
}

// Wikilink represents a wikilink node ([[page]] or [[page|display]])
type Wikilink struct {
	*BaseNode
	Target      string
	DisplayText string
}

// NewWikilink creates a new wikilink node
func NewWikilink(target, displayText string, rng Range) *Wikilink {
	return &Wikilink{
		BaseNode:    NewBaseNode(NodeWikilink, rng),
		Target:      target,
		DisplayText: displayText,
	}
}

// List represents a list node
type List struct {
	*BaseNode
	Ordered bool
	Tight   bool
}

// NewList creates a new list node
func NewList(ordered, tight bool, rng Range) *List {
	return &List{
		BaseNode: NewBaseNode(NodeList, rng),
		Ordered:  ordered,
		Tight:    tight,
	}
}

// ListItem represents a list item node
type ListItem struct {
	*BaseNode
	TaskList bool
	Checked  bool
}

// NewListItem creates a new list item node
func NewListItem(taskList, checked bool, rng Range) *ListItem {
	return &ListItem{
		BaseNode: NewBaseNode(NodeListItem, rng),
		TaskList: taskList,
		Checked:  checked,
	}
}

// Emphasis represents emphasized text (*text* or _text_)
type Emphasis struct {
	*BaseNode
}

// NewEmphasis creates a new emphasis node
func NewEmphasis(rng Range) *Emphasis {
	return &Emphasis{
		BaseNode: NewBaseNode(NodeEmphasis, rng),
	}
}

// Strong represents strong text (**text** or __text__)
type Strong struct {
	*BaseNode
}

// NewStrong creates a new strong node
func NewStrong(rng Range) *Strong {
	return &Strong{
		BaseNode: NewBaseNode(NodeStrong, rng),
	}
}

// Code represents inline code (`code`)
type Code struct {
	*BaseNode
	Content string
}

// NewCode creates a new code node
func NewCode(content string, rng Range) *Code {
	return &Code{
		BaseNode: NewBaseNode(NodeCode, rng),
		Content:  content,
	}
}