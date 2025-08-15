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

// AddChild overrides BaseNode.AddChild to set the concrete Document as parent
func (d *Document) AddChild(child Node) {
	child.SetParent(d)
	d.BaseNode.children = append(d.BaseNode.children, child)
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

// Accept implements the visitor pattern for Heading
func (h *Heading) Accept(visitor Visitor) error {
	return visitor.Visit(h)
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

// Accept implements the visitor pattern for CodeBlock
func (cb *CodeBlock) Accept(visitor Visitor) error {
	return visitor.Visit(cb)
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

// Accept implements the visitor pattern for Wikilink
func (w *Wikilink) Accept(visitor Visitor) error {
	return visitor.Visit(w)
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

// Accept implements the visitor pattern for List
func (l *List) Accept(visitor Visitor) error {
	return visitor.Visit(l)
}

// AddChild overrides BaseNode.AddChild to set the concrete List as parent
func (l *List) AddChild(child Node) {
	child.SetParent(l)
	l.BaseNode.children = append(l.BaseNode.children, child)
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

// Accept implements the visitor pattern for ListItem
func (li *ListItem) Accept(visitor Visitor) error {
	return visitor.Visit(li)
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

// FindListItemAtLine finds the list item that contains the given line number
func (d *Document) FindListItemAtLine(line int) *ListItem {
	var result *ListItem
	
	walker := NewWalker(WalkFunc(func(node Node) error {
		if listItem, ok := node.(*ListItem); ok {
			// Check if this list item's range contains the target line
			if listItem.Range().Start.Line <= line && line <= listItem.Range().End.Line {
				result = listItem
			}
		}
		return nil
	}))
	
	walker.Walk(d)
	return result
}

// FindParentList finds the parent List node for a given ListItem
func (li *ListItem) FindParentList() *List {
	parent := li.Parent()
	for parent != nil {
		// Check if parent is a List by node type
		if parent.Type() == NodeList {
			// Since all our nodes embed BaseNode, we need to walk up to find the concrete List
			// In our case, the parent should be the concrete List node
			if list, ok := parent.(*List); ok {
				return list
			}
		}
		parent = parent.Parent()
	}
	return nil
}

// GetListItems returns all ListItem children of this List
func (l *List) GetListItems() []*ListItem {
	var items []*ListItem
	for _, child := range l.Children() {
		if listItem, ok := child.(*ListItem); ok {
			items = append(items, listItem)
		}
	}
	return items
}

// GetSiblingListItems returns all sibling ListItems in the same parent List
func (li *ListItem) GetSiblingListItems() []*ListItem {
	parentList := li.FindParentList()
	if parentList == nil {
		return nil
	}
	return parentList.GetListItems()
}

// GetListItemIndex returns the index of this ListItem within its parent List
func (li *ListItem) GetListItemIndex() int {
	siblings := li.GetSiblingListItems()
	for i, sibling := range siblings {
		if sibling == li {
			return i
		}
	}
	return -1
}
