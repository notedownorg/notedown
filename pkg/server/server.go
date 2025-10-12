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

	"github.com/notedownorg/notedown/apis/go/application_server/v1alpha1"
	"github.com/notedownorg/notedown/pkg/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DocumentServer implements the DocumentService gRPC interface
type DocumentServer struct {
	v1alpha1.UnimplementedDocumentServiceServer

	workspaceRoot       string
	workspaceDiscoverer *workspaceDiscoverer
	documentLoader      *DocumentLoader
}

// NewDocumentServer creates a new DocumentService server
func NewDocumentServer(workspaceRoot string) (*DocumentServer, error) {
	if workspaceRoot == "" {
		return nil, fmt.Errorf("workspace root cannot be empty")
	}

	// Validate workspace root exists and find actual root if needed
	actualRoot, err := config.FindWorkspaceRoot(workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to find workspace root: %w", err)
	}
	if actualRoot == "" {
		// If no .notedown directory found, use the provided path as-is
		actualRoot = workspaceRoot
	}

	discoverer := newWorkspaceDiscoverer(actualRoot)

	return &DocumentServer{
		workspaceRoot:       actualRoot,
		workspaceDiscoverer: discoverer,
		documentLoader:      NewDocumentLoader(),
	}, nil
}

// ListDocuments implements the ListDocuments RPC method
func (ds *DocumentServer) ListDocuments(ctx context.Context, req *v1alpha1.ListDocumentsRequest) (*v1alpha1.ListDocumentsResponse, error) {
	// Discover all markdown files in workspace via channels
	filesChan, errChan := ds.workspaceDiscoverer.discoverDocuments()

	// Process documents through the fan-out/fan-in pipeline
	documents, err := ds.documentLoader.processDocumentsPipeline(ctx, filesChan, req.Filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process documents: %v", err)
	}

	// Check for any discovery errors that occurred during processing
	select {
	case err := <-errChan:
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to discover documents: %v", err)
		}
	default:
		// No error
	}

	return &v1alpha1.ListDocumentsResponse{Documents: documents}, nil
}

// GetWorkspaceRoot returns the workspace root path
func (ds *DocumentServer) GetWorkspaceRoot() string {
	return ds.workspaceRoot
}
