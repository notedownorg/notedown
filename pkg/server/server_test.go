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
	"path/filepath"
	"testing"

	"github.com/notedownorg/notedown/apis/go/application_server/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestDocumentServer_ListDocuments(t *testing.T) {
	// Get test workspace path
	testWorkspace := filepath.Join("..", "testdata")

	// Create server
	server, err := NewDocumentServer(testWorkspace)
	require.NoError(t, err)
	require.NotNil(t, server)

	ctx := context.Background()

	t.Run("list all documents", func(t *testing.T) {
		req := &v1alpha1.ListDocumentsRequest{}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should find all markdown files
		assert.Len(t, resp.Documents, 3)

		// Check that documents have expected fields
		for _, doc := range resp.Documents {
			assert.NotEmpty(t, doc.Path)
			assert.NotEmpty(t, doc.Checksum)
			// Some documents have metadata, some don't
		}
	})

	t.Run("filter by metadata - status active", func(t *testing.T) {
		filterValue, err := structpb.NewValue("active")
		require.NoError(t, err)

		req := &v1alpha1.ListDocumentsRequest{
			Filter: &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_MetadataFilter{
					MetadataFilter: &v1alpha1.MetadataFilter{
						Field:    "status",
						Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
						Value:    filterValue,
					},
				},
			},
		}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should only find project-notes.md
		assert.Len(t, resp.Documents, 1)
		assert.Equal(t, "project-notes.md", resp.Documents[0].Path)
	})

	t.Run("filter by metadata - high priority", func(t *testing.T) {
		filterValue, err := structpb.NewValue("high")
		require.NoError(t, err)

		req := &v1alpha1.ListDocumentsRequest{
			Filter: &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_MetadataFilter{
					MetadataFilter: &v1alpha1.MetadataFilter{
						Field:    "priority",
						Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
						Value:    filterValue,
					},
				},
			},
		}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should only find project-notes.md
		assert.Len(t, resp.Documents, 1)
		assert.Equal(t, "project-notes.md", resp.Documents[0].Path)
	})

	t.Run("filter by tag contains", func(t *testing.T) {
		filterValue, err := structpb.NewValue("project")
		require.NoError(t, err)

		req := &v1alpha1.ListDocumentsRequest{
			Filter: &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_MetadataFilter{
					MetadataFilter: &v1alpha1.MetadataFilter{
						Field:    "tags",
						Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_CONTAINS,
						Value:    filterValue,
					},
				},
			},
		}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should find project-notes.md (has "project" tag)
		assert.Len(t, resp.Documents, 1)
		assert.Equal(t, "project-notes.md", resp.Documents[0].Path)
	})

	t.Run("filter by metadata exists", func(t *testing.T) {
		req := &v1alpha1.ListDocumentsRequest{
			Filter: &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_MetadataFilter{
					MetadataFilter: &v1alpha1.MetadataFilter{
						Field:    "title",
						Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EXISTS,
					},
				},
			},
		}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should find documents with frontmatter (project-notes.md and api-design.md)
		assert.Len(t, resp.Documents, 2)
	})

	t.Run("AND filter", func(t *testing.T) {
		tagValue, err := structpb.NewValue("project")
		require.NoError(t, err)

		statusValue, err := structpb.NewValue("active")
		require.NoError(t, err)

		req := &v1alpha1.ListDocumentsRequest{
			Filter: &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_AndFilter{
					AndFilter: &v1alpha1.AndFilter{
						Filters: []*v1alpha1.FilterExpression{
							{
								Expression: &v1alpha1.FilterExpression_MetadataFilter{
									MetadataFilter: &v1alpha1.MetadataFilter{
										Field:    "tags",
										Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_CONTAINS,
										Value:    tagValue,
									},
								},
							},
							{
								Expression: &v1alpha1.FilterExpression_MetadataFilter{
									MetadataFilter: &v1alpha1.MetadataFilter{
										Field:    "status",
										Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
										Value:    statusValue,
									},
								},
							},
						},
					},
				},
			},
		}

		resp, err := server.ListDocuments(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should find project-notes.md (has both "project" tag and "active" status)
		assert.Len(t, resp.Documents, 1)
		assert.Equal(t, "project-notes.md", resp.Documents[0].Path)
	})
}

func TestDocumentServer_DocumentContent(t *testing.T) {
	// Get test workspace path
	testWorkspace := filepath.Join("..", "testdata")

	// Create server
	server, err := NewDocumentServer(testWorkspace)
	require.NoError(t, err)

	ctx := context.Background()

	// Get all documents
	req := &v1alpha1.ListDocumentsRequest{}
	resp, err := server.ListDocuments(ctx, req)
	require.NoError(t, err)

	// Find project-notes.md document
	var projectDoc *v1alpha1.Document
	for _, doc := range resp.Documents {
		if doc.Path == "project-notes.md" {
			projectDoc = doc
			break
		}
	}
	require.NotNil(t, projectDoc, "project-notes.md should be found")

	t.Run("metadata extraction", func(t *testing.T) {
		require.NotNil(t, projectDoc.Metadata)

		// Check title
		title := projectDoc.Metadata.Fields["title"].GetStringValue()
		assert.Equal(t, "Project Notes", title)

		// Check status
		status := projectDoc.Metadata.Fields["status"].GetStringValue()
		assert.Equal(t, "active", status)

		// Check priority
		priority := projectDoc.Metadata.Fields["priority"].GetStringValue()
		assert.Equal(t, "high", priority)

		// Check tags array
		tags := projectDoc.Metadata.Fields["tags"].GetListValue()
		require.NotNil(t, tags)
		assert.Len(t, tags.Values, 2)
		assert.Equal(t, "project", tags.Values[0].GetStringValue())
		assert.Equal(t, "planning", tags.Values[1].GetStringValue())
	})

	t.Run("wikilinks extraction", func(t *testing.T) {
		require.Len(t, projectDoc.Wikilinks, 2)

		// Check first wikilink
		wl1 := projectDoc.Wikilinks[0]
		assert.Equal(t, "project-plan", wl1.Target)
		assert.Empty(t, wl1.DisplayText) // No pipe notation
		assert.Greater(t, int(wl1.Line), 0)
		assert.Greater(t, int(wl1.Column), 0)

		// Check second wikilink (with pipe notation)
		wl2 := projectDoc.Wikilinks[1]
		assert.Equal(t, "api-design", wl2.Target)
		assert.Equal(t, "API Design Document", wl2.DisplayText)
		assert.Greater(t, int(wl2.Line), 0)
		assert.Greater(t, int(wl2.Column), 0)
	})

	t.Run("tasks extraction", func(t *testing.T) {
		require.Len(t, projectDoc.Tasks, 3)

		// Check task states
		expectedStates := []string{" ", "x", "wip"}
		for i, task := range projectDoc.Tasks {
			assert.Equal(t, expectedStates[i], task.State)
			assert.NotEmpty(t, task.Text)
			assert.Greater(t, int(task.Line), 0)
			assert.Greater(t, int(task.Column), 0)
		}
	})

	t.Run("checksum generation", func(t *testing.T) {
		// Checksum should be a 64-character hex string (SHA-256)
		assert.Len(t, projectDoc.Checksum, 64)
		assert.Regexp(t, "^[a-f0-9]{64}$", projectDoc.Checksum)
	})
}
