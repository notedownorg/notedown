// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extensions

import (
	"fmt"
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
	hasPipe := false
	var concealStart, concealEnd int

	// Check for pipe separator
	if pipePos := strings.Index(target, "|"); pipePos != -1 {
		hasPipe = true
		displayText = strings.TrimSpace(target[pipePos+1:])
		target = strings.TrimSpace(target[:pipePos])

		// Calculate conceal range: from after [[ to before |
		// concealStart is relative to the start of the wikilink (after [[)
		concealStart = 2         // Start after [[
		concealEnd = 2 + pipePos // End before |
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

	// Create wikilink AST node with segment information
	_, segment := block.PeekLine()
	startOffset := segment.Start
	endOffset := startOffset + wikilinkLength

	node := &WikilinkAST{
		Target:       target,
		DisplayText:  displayText,
		HasPipe:      hasPipe,
		ConcealStart: concealStart,
		ConcealEnd:   concealEnd,
		segment:      text.NewSegment(startOffset, endOffset),
	}

	// Advance the reader by the length of the wikilink
	block.Advance(wikilinkLength)

	return node
}

// WikilinkAST represents a wikilink in the goldmark AST
type WikilinkAST struct {
	ast.BaseInline
	Target       string
	DisplayText  string
	HasPipe      bool         // Whether this wikilink has a pipe separator
	ConcealStart int          // Start position of concealable range (relative to wikilink start)
	ConcealEnd   int          // End position of concealable range (relative to wikilink start)
	segment      text.Segment // Position information
}

// Segment returns the text segment of this wikilink
func (n *WikilinkAST) Segment() text.Segment {
	return n.segment
}

// Dump implements ast.Node
func (n *WikilinkAST) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Target":       n.Target,
		"DisplayText":  n.DisplayText,
		"HasPipe":      fmt.Sprintf("%v", n.HasPipe),
		"ConcealStart": fmt.Sprintf("%d", n.ConcealStart),
		"ConcealEnd":   fmt.Sprintf("%d", n.ConcealEnd),
	}, nil)
}

// Kind returns the node kind
func (n *WikilinkAST) Kind() ast.NodeKind {
	return WikilinkKind
}

// WikilinkKind is the kind for wikilink nodes
var WikilinkKind = ast.NewNodeKind("Wikilink")
