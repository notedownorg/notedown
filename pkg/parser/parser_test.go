package parser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("Expected parser to be created, got nil")
	}
}

func TestParseSimpleMarkdown(t *testing.T) {
	parser := NewParser()
	source := `# Hello World

This is a paragraph with **bold** and *italic* text.

## Subheading

Another paragraph with ` + "`inline code`" + `.`

	doc, err := parser.ParseString(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if doc.Type() != NodeDocument {
		t.Errorf("Expected document node, got %v", doc.Type())
	}

	if len(doc.Children()) == 0 {
		t.Error("Expected document to have children")
	}
	
	// Debug: Print children info
	t.Logf("Document has %d children", len(doc.Children()))
	for i, child := range doc.Children() {
		t.Logf("Child %d: %s", i, child.Type())
	}
}

func TestParseWikilink(t *testing.T) {
	parser := NewParser()
	source := `This paragraph contains a [[wikilink]] and a [[target|display text]] link.`

	doc, err := parser.ParseString(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Find wikilinks by walking the tree
	var wikilinks []*Wikilink
	walker := NewWalker(WalkFunc(func(node Node) error {
		if wikilink, ok := node.(*Wikilink); ok {
			wikilinks = append(wikilinks, wikilink)
		}
		return nil
	}))

	if err := walker.Walk(doc); err != nil {
		t.Fatalf("Error walking tree: %v", err)
	}

	if len(wikilinks) != 2 {
		t.Errorf("Expected 2 wikilinks, got %d", len(wikilinks))
		return
	}

	// Test first wikilink
	if wikilinks[0].Target != "wikilink" {
		t.Errorf("Expected target 'wikilink', got '%s'", wikilinks[0].Target)
	}
	if wikilinks[0].DisplayText != "wikilink" {
		t.Errorf("Expected display text 'wikilink', got '%s'", wikilinks[0].DisplayText)
	}

	// Test second wikilink
	if wikilinks[1].Target != "target" {
		t.Errorf("Expected target 'target', got '%s'", wikilinks[1].Target)
	}
	if wikilinks[1].DisplayText != "display text" {
		t.Errorf("Expected display text 'display text', got '%s'", wikilinks[1].DisplayText)
	}
}

func TestParseCodeBlock(t *testing.T) {
	parser := NewParser()
	source := "```go\npackage main\n\nfunc main() {}\n```"

	doc, err := parser.ParseString(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Find code blocks
	var codeBlocks []*CodeBlock
	walker := NewWalker(WalkFunc(func(node Node) error {
		if block, ok := node.(*CodeBlock); ok {
			codeBlocks = append(codeBlocks, block)
		}
		return nil
	}))

	if err := walker.Walk(doc); err != nil {
		t.Fatalf("Error walking tree: %v", err)
	}

	if len(codeBlocks) != 1 {
		t.Fatalf("Expected 1 code block, got %d", len(codeBlocks))
	}

	block := codeBlocks[0]
	if block.Language != "go" {
		t.Errorf("Expected language 'go', got '%s'", block.Language)
	}
	if !block.Fenced {
		t.Error("Expected fenced code block")
	}
}

func TestParseHeadings(t *testing.T) {
	parser := NewParser()
	source := `# Heading 1
## Heading 2
### Heading 3`

	doc, err := parser.ParseString(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Debug: Print all nodes with more detail
	walker := NewWalker(WalkFunc(func(node Node) error {
		if heading, ok := node.(*Heading); ok {
			t.Logf("Found node: %s (Heading level %d, text '%s')", node.Type(), heading.Level, heading.Text)
		} else {
			t.Logf("Found node: %s", node.Type())
		}
		return nil
	}))

	if err := walker.Walk(doc); err != nil {
		t.Fatalf("Error walking tree: %v", err)
	}

	var headings []*Heading
	walker = NewWalker(WalkFunc(func(node Node) error {
		t.Logf("Checking node type %s for heading (concrete type: %T)", node.Type(), node)
		if node.Type() == NodeHeading {
			if heading, ok := node.(*Heading); ok {
				headings = append(headings, heading)
				t.Logf("Found heading: level %d, text '%s'", heading.Level, heading.Text)
			} else {
				t.Logf("Node has NodeHeading type but failed cast to *Heading (concrete type: %T)", node)
			}
		}
		return nil
	}))

	if err := walker.Walk(doc); err != nil {
		t.Fatalf("Error walking tree: %v", err)
	}

	if len(headings) != 3 {
		t.Fatalf("Expected 3 headings, got %d", len(headings))
	}

	expectedLevels := []int{1, 2, 3}
	expectedTexts := []string{"Heading 1", "Heading 2", "Heading 3"}

	for i, heading := range headings {
		if heading.Level != expectedLevels[i] {
			t.Errorf("Expected heading level %d, got %d", expectedLevels[i], heading.Level)
		}
		if heading.Text != expectedTexts[i] {
			t.Errorf("Expected heading text '%s', got '%s'", expectedTexts[i], heading.Text)
		}
	}
}

func TestNodePositions(t *testing.T) {
	parser := NewParser()
	source := `# Title

Paragraph`

	doc, err := parser.ParseString(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check document range
	docRange := doc.Range()
	if docRange.Start.Line != 1 || docRange.Start.Column != 1 {
		t.Errorf("Expected document start at line 1, column 1, got line %d, column %d", 
			docRange.Start.Line, docRange.Start.Column)
	}

	// Find the first heading
	var firstHeading *Heading
	for _, child := range doc.Children() {
		if heading, ok := child.(*Heading); ok {
			firstHeading = heading
			break
		}
	}

	if firstHeading == nil {
		t.Fatal("Expected to find a heading")
	}

	headingRange := firstHeading.Range()
	if headingRange.Start.Line != 1 {
		t.Errorf("Expected heading on line 1, got line %d", headingRange.Start.Line)
	}
}