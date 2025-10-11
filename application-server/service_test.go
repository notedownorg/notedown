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
	"os"
	"path/filepath"
	"testing"

	pb "github.com/notedownorg/notedown/notedown/application_server/v1"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestDocumentService(t *testing.T) {
	// Create a temporary workspace with test documents
	tempDir := t.TempDir()

	// Create test documents with frontmatter
	testDocs := map[string]string{
		"README.md": `---
title: "README"
author: "Test Author"
tags: ["readme", "documentation"]
priority: 1
published: true
---

# README

This is a test README file.
`,
		"docs/guide.md": `---
title: "User Guide"
author: "Test Author"
tags: ["documentation", "guide"]
priority: 2
published: true
---

# User Guide

This is a user guide.
`,
		"notes/draft.md": `---
title: "Draft Notes"
author: "Test Author"
tags: ["notes", "draft"]
priority: 3
published: false
---

# Draft Notes

These are draft notes.
`,
		"blog/post1.md": `---
title: "Blog Post 1"
author: "Different Author"
tags: ["blog", "post"]
priority: 1
published: true
---

# Blog Post 1

This is a blog post.
`,
	}

	// Create the test files
	for path, content := range testDocs {
		fullPath := filepath.Join(tempDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Create the document service
	logger := log.New(os.Stderr, log.Error) // Use error level to reduce noise
	service := NewDocumentService(logger, []string{tempDir})

	t.Run("list all documents", func(t *testing.T) {
		req := &pb.ListDocumentsRequest{}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Len(t, resp.Documents, 4)
	})

	t.Run("filter by metadata - equals", func(t *testing.T) {
		filter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Test Author")
		require.NoError(t, err)

		req := &pb.ListDocumentsRequest{
			Filter: filter,
		}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)

		assert.Len(t, resp.Documents, 3)
	})

	t.Run("filter by metadata - contains tag", func(t *testing.T) {
		filter, err := NewMetadataFilter("tags", pb.MetadataOperator_METADATA_OPERATOR_CONTAINS, "documentation")
		require.NoError(t, err)

		req := &pb.ListDocumentsRequest{
			Filter: filter,
		}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)

		// README.md and docs/guide.md
		assert.Len(t, resp.Documents, 2)
	})

	t.Run("complex AND filter", func(t *testing.T) {
		authorFilter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Test Author")
		require.NoError(t, err)

		publishedFilter, err := NewMetadataFilter("published", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, true)
		require.NoError(t, err)

		andFilter := NewAndFilter(authorFilter, publishedFilter)

		req := &pb.ListDocumentsRequest{
			Filter: andFilter,
		}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)

		// README.md and docs/guide.md (not draft.md which is unpublished)
		assert.Len(t, resp.Documents, 2)
	})

	t.Run("complex OR filter", func(t *testing.T) {
		priorityHighFilter, err := NewMetadataFilter("priority", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, 1)
		require.NoError(t, err)

		authorDifferentFilter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Different Author")
		require.NoError(t, err)

		orFilter := NewOrFilter(priorityHighFilter, authorDifferentFilter)

		req := &pb.ListDocumentsRequest{
			Filter: orFilter,
		}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)

		// Should match: README.md (priority 1) and blog/post1.md (priority 1 AND different author)

		assert.Len(t, resp.Documents, 2)

		// Debug: Print what we actually got
		t.Logf("Got %d documents:", len(resp.Documents))
		for _, doc := range resp.Documents {
			t.Logf("- %s", doc.Path)
		}
	})

	t.Run("NOT filter", func(t *testing.T) {
		publishedFilter, err := NewMetadataFilter("published", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, true)
		require.NoError(t, err)

		notFilter := NewNotFilter(publishedFilter)

		req := &pb.ListDocumentsRequest{
			Filter: notFilter,
		}

		resp, err := service.ListDocuments(context.Background(), req)
		require.NoError(t, err)

		assert.Len(t, resp.Documents, 1)
		assert.Contains(t, resp.Documents[0].Path, "notes/draft.md")
	})

}

func TestFilterEngine(t *testing.T) {
	filterEngine := newFilterEngine()

	// Create test document
	metadata, err := structpb.NewStruct(map[string]any{
		"title":     "Test Document",
		"author":    "Test Author",
		"tags":      []any{"test", "example"},
		"priority":  5,
		"published": true,
	})
	require.NoError(t, err)

	doc := &pb.Document{
		Path:     "test.md",
		Metadata: metadata,
	}

	t.Run("metadata exists", func(t *testing.T) {
		filter, err := NewMetadataFilter("title", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(filterEngine, []*pb.Document{doc}, filter)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
	})

	t.Run("metadata not exists", func(t *testing.T) {
		filter, err := NewMetadataFilter("nonexistent", pb.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS, nil)
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(filterEngine, []*pb.Document{doc}, filter)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
	})

	t.Run("string contains", func(t *testing.T) {
		filter, err := NewMetadataFilter("title", pb.MetadataOperator_METADATA_OPERATOR_CONTAINS, "Test")
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(filterEngine, []*pb.Document{doc}, filter)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
	})

	t.Run("numeric comparison", func(t *testing.T) {
		filter, err := NewMetadataFilter("priority", pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN, 3)
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(filterEngine, []*pb.Document{doc}, filter)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
	})

	t.Run("array contains", func(t *testing.T) {
		filter, err := NewMetadataFilter("tags", pb.MetadataOperator_METADATA_OPERATOR_CONTAINS, "test")
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(filterEngine, []*pb.Document{doc}, filter)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
	})
}
