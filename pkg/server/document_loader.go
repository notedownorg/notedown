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

package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/notedownorg/notedown/apis/go/application_server/v1alpha1"
	"github.com/notedownorg/notedown/pkg/parser"
	"google.golang.org/protobuf/types/known/structpb"
)

// ProcessedDocument represents a fully processed document with extracted content
type ProcessedDocument struct {
	Path      string
	Checksum  string
	Metadata  map[string]any
	Wikilinks []*v1alpha1.Wikilink
	Tasks     []*v1alpha1.Task
	Error     error
}

// ParsedDocument represents a document that has been parsed but not fully processed
type ParsedDocument struct {
	Path     string
	Checksum string
	Metadata map[string]any
	Document *parser.Document
	Error    error
}

// DocumentLoader handles parallel loading and processing of documents
type DocumentLoader struct {
	parser parser.Parser
}

// NewDocumentLoader creates a new document loader
func NewDocumentLoader() *DocumentLoader {
	return &DocumentLoader{
		parser: parser.NewParser(),
	}
}

// ProcessDocumentsPipeline processes documents using a fan-out/fan-in pipeline
func (dl *DocumentLoader) processDocumentsPipeline(ctx context.Context, filesChan <-chan *DocumentFile, filter *v1alpha1.FilterExpression) ([]*v1alpha1.Document, error) {
	parsedChan := make(chan *ParsedDocument)
	filteredChan := make(chan *ParsedDocument)
	resultsChan := make(chan *v1alpha1.Document)

	// Start pipeline stages with proper coordination
	var parseWG, filterWG, extractWG sync.WaitGroup

	// Stage 1 → 2: Parse documents (fan-out)
	numParsers := 20 // Cap parallelism for very large workspaces
	for range numParsers {
		parseWG.Add(1)
		go dl.parseStage(ctx, filesChan, parsedChan, &parseWG)
	}

	// Stage 2 → 3: Filter documents (single goroutine)
	filterWG.Add(1)
	go dl.filterStage(ctx, parsedChan, filteredChan, filter, &filterWG)

	// Stage 3 → 4: Extract content (fan-out)
	numExtractors := 10 // Reasonable parallelism for content extraction
	for i := 0; i < numExtractors; i++ {
		extractWG.Add(1)
		go dl.extractStage(ctx, filteredChan, resultsChan, &extractWG)
	}

	// Coordinate pipeline stage shutdown
	go func() {
		// Close parsedChan when all parsers are done
		parseWG.Wait()
		close(parsedChan)
	}()

	go func() {
		// filteredChan is closed by filterStage when parsedChan closes
		filterWG.Wait()
	}()

	go func() {
		// Close resultsChan when all extractors are done
		extractWG.Wait()
		close(resultsChan)
	}()

	var results []*v1alpha1.Document
	for result := range resultsChan {
		results = append(results, result)
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return results, nil
}

// ParseDocuments parses documents without extracting wikilinks/tasks (for filtering)
func (dl *DocumentLoader) ParseDocuments(ctx context.Context, files []*DocumentFile) ([]*ParsedDocument, error) {
	if len(files) == 0 {
		return nil, nil
	}

	// Create unbuffered channels for fan-out/fan-in pattern
	filesChan := make(chan *DocumentFile)
	resultsChan := make(chan *ParsedDocument)

	// Start one goroutine per file (fan-out)
	var wg sync.WaitGroup

	// Launch workers - one per file for maximum parallelism
	for range files {
		wg.Add(1)
		go dl.parseWorker(ctx, filesChan, resultsChan, &wg)
	}

	// Send files to workers
	go func() {
		defer close(filesChan)
		for _, file := range files {
			select {
			case filesChan <- file:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect results (fan-in)
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Gather all results
	var results []*ParsedDocument
	for result := range resultsChan {
		results = append(results, result)
	}

	// Check for context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return results, nil
}

// parseWorker processes a single file for parsing only (no wikilink/task extraction)
func (dl *DocumentLoader) parseWorker(ctx context.Context, filesChan <-chan *DocumentFile, resultsChan chan<- *ParsedDocument, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case file, ok := <-filesChan:
		if !ok {
			return // Channel closed
		}

		result := dl.parseFileOnly(file)

		select {
		case resultsChan <- result:
		case <-ctx.Done():
			return
		}

	case <-ctx.Done():
		return
	}
}

// parseFileOnly parses a file and extracts metadata only (for filtering)
func (dl *DocumentLoader) parseFileOnly(file *DocumentFile) *ParsedDocument {
	result := &ParsedDocument{
		Path:     file.Path,
		Checksum: file.Checksum,
	}

	// Read file content
	content, err := dl.readFile(file.AbsPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to read file %s: %w", file.Path, err)
		return result
	}

	// Parse document
	doc, err := dl.parser.Parse(content)
	if err != nil {
		result.Error = fmt.Errorf("failed to parse file %s: %w", file.Path, err)
		return result
	}

	// Store metadata and parsed document
	result.Metadata = doc.Metadata
	result.Document = doc

	return result
}

// ExtractContent extracts wikilinks and tasks from a parsed document
func (dl *DocumentLoader) ExtractContent(parsed *ParsedDocument) *ProcessedDocument {
	result := &ProcessedDocument{
		Path:     parsed.Path,
		Checksum: parsed.Checksum,
		Metadata: parsed.Metadata,
		Error:    parsed.Error,
	}

	if parsed.Error != nil {
		return result
	}

	// Extract wikilinks and tasks
	result.Wikilinks = dl.extractWikilinks(parsed.Document)
	result.Tasks = dl.extractTasks(parsed.Document)

	return result
}

// parseStage processes files through the parsing stage of the pipeline
func (dl *DocumentLoader) parseStage(ctx context.Context, filesChan <-chan *DocumentFile, parsedChan chan<- *ParsedDocument, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case file, ok := <-filesChan:
			if !ok {
				return // Input channel closed
			}

			parsed := dl.parseFileOnly(file)

			select {
			case parsedChan <- parsed:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// filterStage filters documents based on metadata
func (dl *DocumentLoader) filterStage(ctx context.Context, parsedChan <-chan *ParsedDocument, filteredChan chan<- *ParsedDocument, filter *v1alpha1.FilterExpression, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(filteredChan) // Close output when done

	for {
		select {
		case parsed, ok := <-parsedChan:
			if !ok {
				return // Input channel closed
			}

			// Skip documents with errors
			if parsed.Error != nil {
				continue
			}

			// Apply filter if provided
			if filter != nil {
				matches, err := EvaluateFilter(filter, parsed.Metadata)
				if err != nil || !matches {
					continue // Skip filtered out documents
				}
			}

			// Document passes filter, send to next stage
			select {
			case filteredChan <- parsed:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// extractStage extracts wikilinks and tasks from filtered documents and converts to protobuf
func (dl *DocumentLoader) extractStage(ctx context.Context, filteredChan <-chan *ParsedDocument, resultsChan chan<- *v1alpha1.Document, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case parsed, ok := <-filteredChan:
			if !ok {
				return // Input channel closed
			}

			processed := dl.ExtractContent(parsed)
			protoDoc, err := dl.toProtoDocument(processed)
			if err != nil {
				// Skip documents that can't be converted
				continue
			}

			select {
			case resultsChan <- protoDoc:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// readFile reads file content
func (dl *DocumentLoader) readFile(path string) ([]byte, error) {
	// #nosec G304 - path is from trusted workspace discovery
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't fail the operation
			_ = err // Explicitly ignore error
		}
	}()

	return io.ReadAll(file)
}

// extractWikilinks extracts all wikilinks from the parsed document
func (dl *DocumentLoader) extractWikilinks(doc *parser.Document) []*v1alpha1.Wikilink {
	var wikilinks []*v1alpha1.Wikilink

	walker := parser.NewWalker(parser.WalkFunc(func(node parser.Node) error {
		if wikilink, ok := node.(*parser.Wikilink); ok {
			// Only set DisplayText if pipe notation was used
			displayText := ""
			if wikilink.HasPipe {
				displayText = wikilink.DisplayText
			}

			// Safe conversion with bounds checking
			line := wikilink.Range().Start.Line
			column := wikilink.Range().Start.Column
			if line > 2147483647 {
				line = 2147483647
			}
			if column > 2147483647 {
				column = 2147483647
			}

			wikilinks = append(wikilinks, &v1alpha1.Wikilink{
				Target:      wikilink.Target,
				DisplayText: displayText,
				Line:        int32(line),   // #nosec G115 - bounds checked above
				Column:      int32(column), // #nosec G115 - bounds checked above
			})
		}
		return nil
	}))

	_ = walker.Walk(doc)
	return wikilinks
}

// extractTasks extracts all tasks from the parsed document
func (dl *DocumentLoader) extractTasks(doc *parser.Document) []*v1alpha1.Task {
	var tasks []*v1alpha1.Task

	walker := parser.NewWalker(parser.WalkFunc(func(node parser.Node) error {
		if listItem, ok := node.(*parser.ListItem); ok && listItem.TaskList {
			// Extract task text from the list item's children
			taskText := dl.extractTaskText(listItem)

			// Safe conversion with bounds checking
			line := listItem.Range().Start.Line
			column := listItem.Range().Start.Column
			if line > 2147483647 {
				line = 2147483647
			}
			if column > 2147483647 {
				column = 2147483647
			}

			tasks = append(tasks, &v1alpha1.Task{
				State:  listItem.TaskState,
				Text:   taskText,
				Line:   int32(line),   // #nosec G115 - bounds checked above
				Column: int32(column), // #nosec G115 - bounds checked above
			})
		}
		return nil
	}))

	_ = walker.Walk(doc)
	return tasks
}

// extractTaskText extracts the text content from a task list item
func (dl *DocumentLoader) extractTaskText(listItem *parser.ListItem) string {
	var text string

	walker := parser.NewWalker(parser.WalkFunc(func(node parser.Node) error {
		if textNode, ok := node.(*parser.Text); ok {
			text += textNode.Content
		}
		return nil
	}))

	_ = walker.Walk(listItem)
	return text
}

// toProtoDocument converts a ProcessedDocument to protobuf Document
func (dl *DocumentLoader) toProtoDocument(doc *ProcessedDocument) (*v1alpha1.Document, error) {
	// Convert metadata to protobuf Struct
	var metadata *structpb.Struct
	if len(doc.Metadata) > 0 {
		// Filter out unsupported types like time.Time
		filteredMetadata := dl.filterSupportedMetadata(doc.Metadata)
		if len(filteredMetadata) > 0 {
			var err error
			metadata, err = structpb.NewStruct(filteredMetadata)
			if err != nil {
				return nil, fmt.Errorf("failed to convert metadata to protobuf: %w", err)
			}
		}
	}

	return &v1alpha1.Document{
		Path:      doc.Path,
		Checksum:  doc.Checksum,
		Metadata:  metadata,
		Wikilinks: doc.Wikilinks,
		Tasks:     doc.Tasks,
	}, nil
}

// filterSupportedMetadata filters out types not supported by protobuf
func (dl *DocumentLoader) filterSupportedMetadata(metadata map[string]any) map[string]any {
	filtered := make(map[string]any)
	for key, value := range metadata {
		if dl.isSupportedType(value) {
			filtered[key] = value
		}
	}
	return filtered
}

// isSupportedType checks if the type is supported by protobuf structpb
func (dl *DocumentLoader) isSupportedType(value any) bool {
	switch value := value.(type) {
	case nil, bool, int, int32, int64, float32, float64, string:
		return true
	case []any:
		// Check if all elements in slice are supported
		slice := value
		for _, item := range slice {
			if !dl.isSupportedType(item) {
				return false
			}
		}
		return true
	case map[string]any:
		// Check if all values in map are supported
		m := value
		for _, v := range m {
			if !dl.isSupportedType(v) {
				return false
			}
		}
		return true
	default:
		// Unsupported types like time.Time
		return false
	}
}
