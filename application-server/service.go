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

package appserver

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"sync"

	pb "github.com/notedownorg/notedown/notedown/application_server/v1"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/notedownorg/notedown/pkg/parser"
	"github.com/notedownorg/notedown/pkg/workspace"
	"google.golang.org/protobuf/types/known/structpb"
)

// DocumentService implements the gRPC DocumentService
type DocumentService struct {
	pb.UnimplementedDocumentServiceServer
	workspaceRoots []string
	workspace      *workspace.Manager
	parser         parser.Parser
	filter         *filterEngine
	logger         *log.Logger
}

// NewDocumentService creates a new DocumentService
func NewDocumentService(logger *log.Logger, workspaceRoots []string) *DocumentService {
	ws := workspace.NewManager(logger)

	// Add workspace roots during initialization
	for _, root := range workspaceRoots {
		if err := ws.AddRoot(root); err != nil {
			logger.Warn("failed to add workspace root", "root", root, "error", err)
			continue
		}
	}

	return &DocumentService{
		workspaceRoots: workspaceRoots,
		workspace:      ws,
		parser:         parser.NewParser(),
		filter:         newFilterEngine(),
		logger:         logger,
	}
}

// ListDocuments implements the ListDocuments RPC method
func (ds *DocumentService) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	ds.logger.Debug("received ListDocuments request")

	// Discover and read all documents from configured workspace roots
	docChan, docErrChan := ds.discoverAndReadDocuments()

	// Filter documents in parallel
	filteredChan, filterErrChan := ds.filter.FilterDocuments(docChan, req.Filter)

	// Collect filtered documents
	var filteredDocs []*pb.Document
	for doc := range filteredChan {
		filteredDocs = append(filteredDocs, doc)
	}

	// Check for any errors that occurred during document discovery
	select {
	case err := <-docErrChan:
		if err != nil {
			return nil, fmt.Errorf("failed to discover documents: %w", err)
		}
	default:
		// No error
	}

	// Check for any errors that occurred during filtering
	select {
	case err := <-filterErrChan:
		if err != nil {
			return nil, fmt.Errorf("failed to filter documents: %w", err)
		}
	default:
		// No error
	}

	ds.logger.Info("ListDocuments completed", "returned", len(filteredDocs))

	return &pb.ListDocumentsResponse{
		Documents: filteredDocs,
	}, nil
}

// discoverAndReadDocuments discovers and reads all markdown files from workspace roots in parallel
// Returns a channel of documents and an error channel
func (ds *DocumentService) discoverAndReadDocuments() (<-chan *pb.Document, <-chan error) {
	docChan := make(chan *pb.Document)
	errChan := make(chan error)

	go func() {
		defer close(docChan)
		defer close(errChan)

		// Discover all markdown files (workspace roots already configured in constructor)
		if err := ds.workspace.DiscoverMarkdownFiles(); err != nil {
			errChan <- err
			return
		}

		files := ds.workspace.GetMarkdownFiles()
		if len(files) == 0 {
			return // No files to process
		}

		// Process files in parallel with unlimited concurrency
		// Each file gets its own goroutine for maximum parallelism
		var wg sync.WaitGroup
		for _, fileInfo := range files {
			wg.Add(1)
			go func(fileInfo *workspace.FileInfo) {
				defer wg.Done()
				doc, err := ds.readAndParseDocument(fileInfo)
				if err != nil {
					ds.logger.Warn("failed to read document", "path", fileInfo.Path, "error", err)
					return
				}
				docChan <- doc
			}(fileInfo)
		}

		// Wait for all goroutines to complete
		wg.Wait()
	}()

	return docChan, errChan
}

// readAndParseDocument reads and parses a single document file
func (ds *DocumentService) readAndParseDocument(fileInfo *workspace.FileInfo) (*pb.Document, error) {
	// Convert URI to file path for reading
	filePath := strings.TrimPrefix(fileInfo.URI, "file://")

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse document to extract metadata
	parsedDoc, err := ds.parser.Parse(content)
	if err != nil {
		return nil, err
	}

	// Calculate checksum
	hash := sha256.Sum256(content)
	checksum := fmt.Sprintf("%x", hash)

	// Extract wikilinks and tasks
	wikilinks := parser.ExtractWikilinks(parsedDoc)
	// TODO: Extract tasks - requires goldmark AST which isn't exposed by parser
	// For now, return empty slice
	var tasks []parser.TaskInfo

	// Convert metadata to protobuf Struct
	metadataStruct, err := structpb.NewStruct(parsedDoc.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metadata to protobuf struct: %w", err)
	}

	// Convert wikilinks to protobuf format
	pbWikilinks := make([]*pb.Wikilink, len(wikilinks))
	for i, wl := range wikilinks {
		pbWikilinks[i] = &pb.Wikilink{
			Target:      wl.Target,
			DisplayText: wl.DisplayText,
			Line:        int32(wl.Line),
			Column:      int32(wl.Column),
		}
	}

	// Convert tasks to protobuf format
	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = &pb.Task{
			State:  task.State,
			Text:   task.Text,
			Line:   int32(task.Line),
			Column: int32(task.Column),
		}
	}

	return &pb.Document{
		Path:      fileInfo.Path,
		Checksum:  checksum,
		Metadata:  metadataStruct,
		Wikilinks: pbWikilinks,
		Tasks:     pbTasks,
	}, nil
}

// Helper functions for building filters programmatically

// NewMetadataFilter creates a new metadata filter
func NewMetadataFilter(field string, operator pb.MetadataOperator, value any) (*pb.FilterExpression, error) {
	// Validate operator
	switch operator {
	case pb.MetadataOperator_METADATA_OPERATOR_EXISTS,
		pb.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS,
		pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
		pb.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS,
		pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
		pb.MetadataOperator_METADATA_OPERATOR_STARTS_WITH,
		pb.MetadataOperator_METADATA_OPERATOR_ENDS_WITH,
		pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
		pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL,
		pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN,
		pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL,
		pb.MetadataOperator_METADATA_OPERATOR_IN,
		pb.MetadataOperator_METADATA_OPERATOR_NOT_IN:
		// Valid operators
	default:
		return nil, fmt.Errorf("unsupported metadata operator: %v", operator)
	}

	pbValue, err := structpb.NewValue(value)
	if err != nil {
		return nil, err
	}

	return &pb.FilterExpression{
		Expression: &pb.FilterExpression_MetadataFilter{
			MetadataFilter: &pb.MetadataFilter{
				Field:    field,
				Operator: operator,
				Value:    pbValue,
			},
		},
	}, nil
}

// NewAndFilter creates a new AND filter
func NewAndFilter(filters ...*pb.FilterExpression) *pb.FilterExpression {
	return &pb.FilterExpression{
		Expression: &pb.FilterExpression_AndFilter{
			AndFilter: &pb.AndFilter{
				Filters: filters,
			},
		},
	}
}

// NewOrFilter creates a new OR filter
func NewOrFilter(filters ...*pb.FilterExpression) *pb.FilterExpression {
	return &pb.FilterExpression{
		Expression: &pb.FilterExpression_OrFilter{
			OrFilter: &pb.OrFilter{
				Filters: filters,
			},
		},
	}
}

// NewNotFilter creates a new NOT filter
func NewNotFilter(filter *pb.FilterExpression) *pb.FilterExpression {
	return &pb.FilterExpression{
		Expression: &pb.FilterExpression_NotFilter{
			NotFilter: &pb.NotFilter{
				Filter: filter,
			},
		},
	}
}
